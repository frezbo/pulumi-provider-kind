// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/frezbo/pulumi-provider-kind/provider/v3/pkg/logging"
	"github.com/frezbo/pulumi-provider-kind/provider/v3/pkg/metadata"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	rpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	pulumilog "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cluster"
)

const (
	kindDefaultProvider = "docker"
	kindDockerProvider  = "docker"
	kindPodmanProvider  = "podman"
)

type cancellationContext struct {
	context context.Context
	cancel  context.CancelFunc
}

func makeCancellationContext() *cancellationContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &cancellationContext{
		context: ctx,
		cancel:  cancel,
	}
}

type kindProvider struct {
	host     *provider.HostClient
	canceler *cancellationContext
	name     string
	version  string
	opts     kindCreateOpts
}

// Partial create options from https://pkg.go.dev/sigs.k8s.io/kind@v0.11.1/pkg/cluster
// Using v1alpha4 schema as default
type kindCreateOpts struct {
	ConfigFile           string
	NodeImage            string
	RetainNodesOnFailure bool
	WaitForNodeReady     time.Duration
	KubeconfigFile       string
	StopBeforeSettingK8s bool
	Provider             string
}

func makeKindProvider(host *provider.HostClient, name, version string) (rpc.ResourceProviderServer, error) {
	// Return the new provider
	return &kindProvider{
		host:     host,
		canceler: makeCancellationContext(),
		name:     name,
		version:  version,
		opts:     kindCreateOpts{},
	}, nil
}

// added as part of pulumi sdk upgrade to v3.6
// not sure what Call does
func (k *kindProvider) Call(ctx context.Context, call *rpc.CallRequest) (*rpc.CallResponse, error) {
	return nil, nil
}

// CheckConfig validates the configuration for this provider.
func (k *kindProvider) CheckConfig(ctx context.Context, req *rpc.CheckRequest) (*rpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.CheckConfig(%s)", k.name, urn)
	pulumilog.V(9).Infof("%s executing", label)

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label:          fmt.Sprintf("%s.news", label),
		RejectUnknowns: true,
		SkipNulls:      true,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "CheckConfig failed because of malformed resource inputs")
	}
	truthyValue := func(argName resource.PropertyKey, props resource.PropertyMap) bool {
		if arg := news[argName]; arg.HasValue() {
			switch {
			case arg.IsString() && len(arg.StringValue()) > 0:
				return true
			default:
				return false
			}
		}
		return false
	}

	var failures []*rpc.CheckFailure

	configFileValueSet := truthyValue("configFile", news)
	if configFileValueSet {
		configFilePath := news["configFile"].StringValue()
		if _, err := os.Stat(configFilePath); err == os.ErrNotExist {
			failures = append(failures, &rpc.CheckFailure{
				Property: "configFile",
				Reason:   fmt.Sprintf("KIND config file does not exist at path: %s", configFilePath),
			})
		}
	}

	providerValueSet := truthyValue("provider", news)

	if providerValueSet {
		switch news["provider"].StringValue() {
		case kindDockerProvider:
		case kindPodmanProvider:
		default:
			failures = append(failures, &rpc.CheckFailure{
				Property: "provider",
				Reason:   "Valid provider values are docker/podman",
			})
		}
	}

	return &rpc.CheckResponse{Inputs: req.GetNews(), Failures: failures}, nil
}

// DiffConfig diffs the configuration for this provider.
func (k *kindProvider) DiffConfig(ctx context.Context, req *rpc.DiffRequest) (*rpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.DiffConfig(%s)", k.name, urn)
	pulumilog.V(9).Infof("%s executing", label)

	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.olds", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, err
	}
	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		RejectAssets: true,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "DiffConfig failed because of malformed resource inputs")
	}

	// any provider changes creates a new kind cluster
	if d := olds.Diff(news); d != nil {
		return &rpc.DiffResponse{
			Changes: rpc.DiffResponse_DIFF_SOME,

			Replaces:            []string{""},
			DeleteBeforeReplace: true,
		}, nil
	}

	return &rpc.DiffResponse{
		Changes: rpc.DiffResponse_DIFF_NONE,
	}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (k *kindProvider) Configure(_ context.Context, req *rpc.ConfigureRequest) (*rpc.ConfigureResponse, error) {
	const trueStr = "true"
	vars := req.GetVariables()

	// set some defaults so that the KIND clusters can be created with any option
	if configFile, exists := vars["kind:config:configFile"]; exists {
		k.opts.ConfigFile = configFile
	}
	if kubeconfig, exists := vars["kind:config:kubeconfigFile"]; exists {
		k.opts.KubeconfigFile = kubeconfig
	}
	if nodeImage, exists := vars["kind:config:nodeImage"]; exists {
		k.opts.NodeImage = nodeImage
	}
	if provider, exists := vars["kind:config:provider"]; exists {
		k.opts.Provider = provider
	} else {
		k.opts.Provider = kindDefaultProvider
	}
	if retain, exists := vars["kind:config:retainNodesOnFailure"]; exists {
		k.opts.RetainNodesOnFailure = retain == trueStr
	} else {
		k.opts.RetainNodesOnFailure = false
	}
	if stopBeforeSettingK8s, exists := vars["kind:config:stopBeforeSettingK8s"]; exists {
		k.opts.StopBeforeSettingK8s = stopBeforeSettingK8s == trueStr
	} else {
		k.opts.StopBeforeSettingK8s = false
	}
	if waitForNodeReady, exists := vars["kind:config:waitForNodeReady"]; exists {
		waitDuration, _ := strconv.Atoi(waitForNodeReady)
		k.opts.WaitForNodeReady = time.Duration(waitDuration) * time.Second
	} else {
		k.opts.WaitForNodeReady = 0 * time.Second
	}

	return &rpc.ConfigureResponse{
		SupportsPreview: true,
	}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (k *kindProvider) Invoke(_ context.Context, req *rpc.InvokeRequest) (*rpc.InvokeResponse, error) {
	tok := req.GetTok()
	return nil, fmt.Errorf("unknown Invoke token '%s'", tok)
}

// StreamInvoke dynamically executes a built-in function in the provider. The result is streamed
// back as a series of messages.
func (k *kindProvider) StreamInvoke(req *rpc.InvokeRequest, server rpc.ResourceProvider_StreamInvokeServer) error {
	tok := req.GetTok()
	return fmt.Errorf("unknown StreamInvoke token '%s'", tok)
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (k *kindProvider) Check(ctx context.Context, req *rpc.CheckRequest) (*rpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.DiffConfig(%s)", k.name, urn)
	pulumilog.V(9).Infof("%s executing", label)

	// Obtain old resource inputs. This is the old version of the resource(s) supplied by the user as
	// an update.
	oldResInputs := req.GetOlds()
	olds, err := plugin.UnmarshalProperties(oldResInputs, plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.olds", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "check failed because malformed resource inputs")
	}

	// old inputs means the resource has been created at-least once
	// so we can assume the input KIND cluster config is actually valid
	oldInputs, err := propMapToKindClusterConfig(olds.Mappable())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert old inputs to kind config")
	}

	// Obtain new resource inputs. This is the new version of the resource(s) supplied by the user as
	// an update.
	newResInputs := req.GetNews()
	news, err := plugin.UnmarshalProperties(newResInputs, plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "check failed because malformed resource inputs")
	}

	// we can't call `propMapToKindClusterConfig()` since we don't know if any inputs
	// is a computed value, let's defer to validation when the actual resource is being created
	// and only add `name` for the resource if it's missing
	newInputs := news.Mappable()

	// Adopt name from old object if appropriate.
	//
	// If the user HAS NOT assigned a name in the new inputs, we autoname it and mark the object as
	// autonamed. This makes it easier for `Diff` to decide whether this
	// needs to be `DeleteBeforeReplace`'d. If the resource is marked `DeleteBeforeReplace`, then
	// `Create` will allocate it a new name later.
	// if the name is not set, we can assume the cluster has not been created
	if oldInputs.Name != "" {
		metadata.AdoptOldAutonameIfUnnamed(newInputs, oldInputs)

	} else {
		metadata.AssignNameIfAutonamable(newInputs, news, urn.Name())
	}

	checkedInputs := resource.NewPropertyMapFromMap(newInputs)

	autonamedInputs, err := plugin.MarshalProperties(checkedInputs, plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, err
	}

	return &rpc.CheckResponse{Inputs: autonamedInputs, Failures: nil}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (k *kindProvider) Diff(ctx context.Context, req *rpc.DiffRequest) (*rpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Diff(%s)", k.name, urn)
	pulumilog.V(9).Infof("%s executing", label)

	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, err
	}

	oldInputs, err := propMapToKindClusterConfig(olds.Mappable())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert old inputs to kind config")
	}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, err
	}

	newInputs, err := propMapToKindClusterConfig(news.Mappable())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert new inputs to kind config")
	}

	if !reflect.DeepEqual(oldInputs, newInputs) {
		return &rpc.DiffResponse{
			Changes:             rpc.DiffResponse_DIFF_SOME,
			Replaces:            []string{""},
			DeleteBeforeReplace: true,
		}, nil
	}

	return &rpc.DiffResponse{}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
func (k *kindProvider) Create(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Create(%s)", k.name, urn)
	pulumilog.V(9).Infof("%s executing", label)

	newInputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, err
	}

	var kindProviderOption cluster.ProviderOption

	switch k.opts.Provider {
	case kindDockerProvider:
		kindProviderOption = cluster.ProviderWithDocker()
	case kindPodmanProvider:
		kindProviderOption = cluster.ProviderWithPodman()
	}

	newInputsMap := newInputs.Mappable()

	if req.GetPreview() {

		newInputsMap["kubeconfig"] = resource.Computed{}
		newInputsMap["name"] = resource.Computed{}

		outputProperties, err := plugin.MarshalProperties(
			resource.NewPropertyMapFromMap(newInputsMap),
			plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true},
		)
		if err != nil {
			return nil, err
		}
		return &rpc.CreateResponse{Properties: outputProperties}, nil
	}

	logger := logging.NewLogger(k.canceler.context, k.host, urn)

	kindProviderConfig := cluster.NewProvider(kindProviderOption, cluster.ProviderWithLogger(logger))

	var kindClusterCreateOptions []cluster.CreateOption

	if k.opts.ConfigFile != "" {
		kindClusterCreateOptions = append(kindClusterCreateOptions, cluster.CreateWithConfigFile(k.opts.ConfigFile))
	}
	if k.opts.KubeconfigFile != "" {
		kindClusterCreateOptions = append(kindClusterCreateOptions, cluster.CreateWithKubeconfigPath(k.opts.KubeconfigFile))
	}
	if k.opts.NodeImage != "" {
		kindClusterCreateOptions = append(kindClusterCreateOptions, cluster.CreateWithNodeImage(k.opts.NodeImage))
	}
	if k.opts.RetainNodesOnFailure {
		kindClusterCreateOptions = append(kindClusterCreateOptions, cluster.CreateWithRetain(k.opts.RetainNodesOnFailure))
	}
	if k.opts.StopBeforeSettingK8s {
		kindClusterCreateOptions = append(kindClusterCreateOptions, cluster.CreateWithStopBeforeSettingUpKubernetes(k.opts.StopBeforeSettingK8s))
	}

	clusterConfig, err := propMapToKindClusterConfig(newInputsMap)
	if err != nil {
		return nil, errors.Wrapf(err, "Create():")
	}

	clusterName := clusterConfig.Name

	kindClusterCreateOptions = append(kindClusterCreateOptions, cluster.CreateWithV1Alpha4Config(clusterConfig))
	kindClusterCreateOptions = append(kindClusterCreateOptions, cluster.CreateWithWaitForReady(k.opts.WaitForNodeReady))

	if err = kindProviderConfig.Create(clusterName, kindClusterCreateOptions...); err != nil {
		// delete any kind cluster that failed to create like the kind cli, unless explicitly set not to
		if !k.opts.RetainNodesOnFailure {
			// a best effort to delete the cluster that failed to create
			// without checking for errors. It's up to the user to cleanup
			// kind clusters that may have been orphaned due to some serious
			// config issues/weird edge cases
			// nolint:errcheck
			kindProviderConfig.Delete(clusterName, k.opts.KubeconfigFile)
		}
		return nil, err
	}

	kubeconfig := ""
	if !k.opts.StopBeforeSettingK8s {
		kubeconfig, err = kindProviderConfig.KubeConfig(clusterName, false)
		if err != nil {
			return nil, err
		}
	}
	newInputsMap["kubeconfig"] = kubeconfig
	newInputsMap["name"] = clusterName

	outputProperties, err := plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(newInputsMap),
		plugin.MarshalOptions{
			Label:        fmt.Sprintf("%s.news", label),
			KeepUnknowns: true,
			SkipNulls:    true,
		},
	)
	if err != nil {
		return nil, err
	}
	return &rpc.CreateResponse{
		Id:         clusterName,
		Properties: outputProperties,
	}, nil
}

// Read the current live state associated with a resource.
func (k *kindProvider) Read(ctx context.Context, req *rpc.ReadRequest) (*rpc.ReadResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Read(%s)", k.name, urn)
	pulumilog.V(9).Infof("%s executing", label)

	panic("Read not implemented for 'kind:cluster:Cluster'")
}

// Update updates an existing resource with new values.
func (k *kindProvider) Update(ctx context.Context, req *rpc.UpdateRequest) (*rpc.UpdateResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Update(%s)", k.name, urn)
	pulumilog.V(9).Infof("%s executing", label)

	return &rpc.UpdateResponse{}, nil
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (k *kindProvider) Delete(ctx context.Context, req *rpc.DeleteRequest) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Delete(%s)", k.name, urn)
	pulumilog.V(9).Infof("%s executing", label)

	var kindProviderOption cluster.ProviderOption

	switch k.opts.Provider {
	case kindDockerProvider:
		kindProviderOption = cluster.ProviderWithDocker()
	case kindPodmanProvider:
		kindProviderOption = cluster.ProviderWithPodman()
	}

	logger := logging.NewLogger(k.canceler.context, k.host, urn)
	provider := cluster.NewProvider(kindProviderOption, cluster.ProviderWithLogger(logger))
	if err := provider.Delete(req.Id, k.opts.KubeconfigFile); err != nil {
		return &pbempty.Empty{}, err
	}

	return &pbempty.Empty{}, nil
}

// Construct creates a new component resource.
func (k *kindProvider) Construct(_ context.Context, _ *rpc.ConstructRequest) (*rpc.ConstructResponse, error) {
	panic("Construct not implemented")
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (k *kindProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*rpc.PluginInfo, error) {
	return &rpc.PluginInfo{
		Version: k.version,
	}, nil
}

// GetSchema returns the JSON-serialized schema for the provider.
func (k *kindProvider) GetSchema(ctx context.Context, req *rpc.GetSchemaRequest) (*rpc.GetSchemaResponse, error) {
	return &rpc.GetSchemaResponse{}, nil
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
// Operations aborted in this way will return an error (e.g., `Update` and `Create` will either a
// creation error or an initialization error). Since Cancel is advisory and non-blocking, it is up
// to the host to decide how long to wait after Cancel is called before (e.g.)
// hard-closing any gRPC connection.
func (k *kindProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	// TODO
	return &pbempty.Empty{}, nil
}

func propMapToKindClusterConfig(inputs map[string]interface{}) (*v1alpha4.Cluster, error) {
	clusterConfig := &v1alpha4.Cluster{}
	clusterConfigData, err := json.Marshal(inputs)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(clusterConfigData, clusterConfig); err != nil {
		return nil, err
	}
	return clusterConfig, nil
}
