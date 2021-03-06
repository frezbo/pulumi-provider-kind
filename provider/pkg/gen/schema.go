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

// Copied and modified heavily based on https://github.com/pulumi/pulumi-kubernetes/blob/v3.5.0/provider/pkg/gen/schema.go
package gen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alecthomas/jsonschema"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	kindconstants "sigs.k8s.io/kind/pkg/cluster/constants"
)

const (
	kindClusterDefinition = "Cluster"
	kindNodeDefinition    = "Node"
)

func PulumiSchema(swagger *jsonschema.Schema) schema.PackageSpec {
	pkg := schema.PackageSpec{
		Name:        "kind",
		Description: "A Pulumi package for creating and managing KIND clusters,",
		License:     "Apache-2.0",
		Keywords:    []string{"pulumi", "kind"},
		Homepage:    "https://github.com/frezbo/pulumi-provider-kind",
		Repository:  "https://github.com/frezbo/pulumi-provider-kind",

		Config: schema.ConfigSpec{
			Variables: map[string]schema.PropertySpec{
				"configFile": {
					TypeSpec:    schema.TypeSpec{Type: "string"},
					Description: "Kind config file to use. Optional",
				},
				"kubeconfigFile": {
					TypeSpec:    schema.TypeSpec{Type: "string"},
					Description: "File to save generated kubeconfig. Default: not set. Optional",
				},
				"nodeImage": {
					TypeSpec:    schema.TypeSpec{Type: "string"},
					Description: "Node image to use. Optional",
				},
				"provider": {
					TypeSpec:    schema.TypeSpec{Type: "string"},
					Description: "Provider to use. Supports docker/podman. Default: docker. Optional",
				},
				"retainNodesOnFailure": {
					TypeSpec:    schema.TypeSpec{Type: "boolean"},
					Description: "Whether to retain the nodes when creation fails. Needs manual cleanup when set to true Default: false. Optional",
				},
				"stopBeforeSettingK8s": {
					TypeSpec:    schema.TypeSpec{Type: "boolean"},
					Description: "Stop before running kubeadm commands. This would need the user to manually retrieve the Kubeconfig. Default: false. Optional",
				},
				"waitForNodeReady": {
					TypeSpec:    schema.TypeSpec{Type: "integer"},
					Description: "Time in seconds to wait for nodes to become ready. Default: none. Optional",
				},
			},
		},
		Provider: schema.ResourceSpec{
			ObjectTypeSpec: schema.ObjectTypeSpec{
				Description: "The provider type for the kind package.",
				Type:        "object",
			},
			InputProperties: map[string]schema.PropertySpec{
				"configFile": {
					TypeSpec:    schema.TypeSpec{Type: "string"},
					Description: "Kind config file to use. Optional",
				},
				"kubeconfigFile": {
					TypeSpec:    schema.TypeSpec{Type: "string"},
					Description: "File to save generated kubeconfig. Default: not set. Optional",
				},
				"nodeImage": {
					TypeSpec:    schema.TypeSpec{Type: "string"},
					Description: "Node image to use. Optional",
				},
				"provider": {
					TypeSpec:    schema.TypeSpec{Type: "string"},
					Description: "Provider to use. Supports docker/podman. Default: docker. Optional",
				},
				"retainNodesOnFailure": {
					TypeSpec:    schema.TypeSpec{Type: "boolean"},
					Description: "Whether to retain the nodes when creation fails. Needs manual cleanup when set to true Default: false. Optional",
				},
				"stopBeforeSettingK8s": {
					TypeSpec:    schema.TypeSpec{Type: "boolean"},
					Description: "Stop before running kubeadm commands. This would need the user to manually retrieve the Kubeconfig. Default: false. Optional",
				},
				"waitForNodeReady": {
					TypeSpec:    schema.TypeSpec{Type: "integer"},
					Description: "Time in seconds to wait for nodes to become ready. Default: none. Optional",
				},
			},
		},
		Types:     map[string]schema.ComplexTypeSpec{},
		Resources: map[string]schema.ResourceSpec{},
		Functions: map[string]schema.FunctionSpec{},
		Language:  map[string]schema.RawMessage{},
	}

	goImportPath := "github.com/frezbo/pulumi-provider-kind/sdk/v3/go/kind"

	pkgImportAliases := map[string]string{}

	for defintion, definitionProperties := range swagger.Definitions {

		defintionSmallCased := strings.ToLower(defintion)

		tok := fmt.Sprintf("kind:%s:%s", defintionSmallCased, defintion)

		objectSpec := schema.ObjectTypeSpec{
			Description: fmt.Sprintf("KIND %s", defintion),
			Type:        "object",
			Properties:  map[string]schema.PropertySpec{},
			Language:    map[string]schema.RawMessage{},
		}

		resourceSpec := schema.ResourceSpec{
			ObjectTypeSpec:  objectSpec,
			InputProperties: map[string]schema.PropertySpec{},
			RequiredInputs:  definitionProperties.Required,
		}

		pkgImportAliases[fmt.Sprintf("%s/%s", goImportPath, defintionSmallCased)] = defintionSmallCased

		for _, definitionPropertyKey := range definitionProperties.Properties.Keys() {
			if val, ok := definitionProperties.Properties.Get(definitionPropertyKey); ok {
				// the returned value is an interface
				// we need to cast it to jsonSchema.Type to access the fields
				definitionPropertyValue := val.(*jsonschema.Type)
				resourceInputProperty := schema.PropertySpec{
					TypeSpec: schema.TypeSpec{
						Type: definitionPropertyValue.Type,
						Ref:  openAPISpecRefToPulumiRef(definitionPropertyValue.Ref),
					},
				}
				if definitionPropertyValue.Items != nil {
					resourceInputProperty.TypeSpec.Items = &schema.TypeSpec{
						Type: definitionPropertyValue.Items.Type,
						Ref:  openAPISpecRefToPulumiRef(definitionPropertyValue.Items.Ref),
					}
				}

				resourceSpec.InputProperties[definitionPropertyKey] = resourceInputProperty

				typeSpec := schema.ComplexTypeSpec{
					ObjectTypeSpec: schema.ObjectTypeSpec{
						Description: fmt.Sprintf("KIND %s type", defintion),
						Type:        "object",
						Properties:  resourceSpec.InputProperties,
						Language:    map[string]schema.RawMessage{},
					},
				}

				// define outputs for the kind cluster
				if defintion == kindClusterDefinition {
					resourceSpec.Properties["kubeconfig"] = schema.PropertySpec{
						Description: "kubeconfig content",
						TypeSpec: schema.TypeSpec{
							Type: "string",
						},
					}

					resourceSpec.Properties["name"] = schema.PropertySpec{
						Description: "cluster name",
						TypeSpec: schema.TypeSpec{
							Type: "string",
						},
					}

					resourceSpec.Required = []string{
						"kubeconfig",
						"name",
					}

					// let's only expose the kind cluster resource
					pkg.Resources[tok] = resourceSpec
					continue
				}
				if defintion == kindNodeDefinition {
					// let's add some constants for the node role type
					// example copied from: https://github.com/pulumi/pulumi-kubernetes/blob/0072954b2cdf088fc2e336ca4c289929f75ec1a5/provider/pkg/gen/overlays.go
					typeSpec.Properties["role"] = schema.PropertySpec{
						Description: "node role type",
						TypeSpec: schema.TypeSpec{
							OneOf: []schema.TypeSpec{
								{
									Type: "string",
								},
								{
									Type: "string",
									Ref:  "#types/kind:node:RoleType",
								},
							},
						},
					}
					pkg.Types["kind:node:RoleType"] = schema.ComplexTypeSpec{
						ObjectTypeSpec: schema.ObjectTypeSpec{
							Type: "string",
						},
						Enum: []schema.EnumValueSpec{
							{
								Name:        "ControlPlane",
								Value:       kindconstants.ControlPlaneNodeRoleValue,
								Description: "node that hosts Kubernetes control-plane components",
							},
							{
								Name:        "Worker",
								Value:       kindconstants.WorkerNodeRoleValue,
								Description: "node that hosts Kubernetes worker",
							},
							{
								Name:        "LoadBalancer",
								Value:       kindconstants.ExternalLoadBalancerNodeRoleValue,
								Description: "node that hosts an external load balancer",
							},
						},
					}
				}
				pkg.Types[tok] = typeSpec
			}
		}

	}

	pkg.Language["go"] = rawMessage(map[string]interface{}{
		"importBasePath":                 goImportPath,
		"packageImportAliases":           pkgImportAliases,
		"generateResourceContainerTypes": true,
	})

	pkg.Language["nodejs"] = rawMessage(map[string]interface{}{
		"dependencies": map[string]string{
			"@pulumi/pulumi": "^3.0.0",
		},
		"devDependencies": map[string]string{
			"typescript":  "^3.7.0",
			"@types/node": "^15.12.0",
		},
	})

	return pkg
}

func rawMessage(v interface{}) schema.RawMessage {
	bytes, err := json.Marshal(v)
	contract.Assert(err == nil)
	return bytes
}

func openAPISpecRefToPulumiRef(ref string) string {
	if ref != "" {
		// convert openAPI schema references to Pulumi schema reference
		// remove `definitions` and replace by `types`
		// ref: https://www.pulumi.com/docs/guides/pulumi-packages/schema/#
		ref = strings.ReplaceAll(ref, "#/definitions/", "")
		ref = fmt.Sprintf("%s/kind:%s:%s", "#/types", strings.ToLower(ref), ref)
		return ref
	}
	return ref
}
