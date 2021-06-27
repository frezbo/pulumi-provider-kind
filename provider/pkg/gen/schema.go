package gen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alecthomas/jsonschema"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/go/common/util/contract"
)

const (
	kindClusterDefinition = "Cluster"
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
				"provider": {
					TypeSpec:    schema.TypeSpec{Type: "string"},
					Description: "Provider to use. Default: docker",
				},
			},
		},
		Provider: schema.ResourceSpec{
			ObjectTypeSpec: schema.ObjectTypeSpec{
				Description: "The provider type for the kind package.",
				Type:        "object",
			},
			InputProperties: map[string]schema.PropertySpec{
				"provider": {
					TypeSpec:    schema.TypeSpec{Type: "string"},
					Description: "Provider to use. Default: docker",
				},
			},
		},
		Types:     map[string]schema.ComplexTypeSpec{},
		Resources: map[string]schema.ResourceSpec{},
		Functions: map[string]schema.FunctionSpec{},
		Language:  map[string]json.RawMessage{},
	}

	goImportPath := "github.com/frezbo/pulumi-provider-kind/sdk/v3/go/kind"

	pkgImportAliases := map[string]string{}

	for defintion, definitionProperties := range swagger.Definitions {

		tok := fmt.Sprintf("kind:index:%s", defintion)

		objectSpec := schema.ObjectTypeSpec{
			Description: fmt.Sprintf("KIND %s", defintion),
			Type:        "object",
			Properties:  map[string]schema.PropertySpec{},
			Language:    map[string]json.RawMessage{},
		}

		resourceSpec := schema.ResourceSpec{
			ObjectTypeSpec:  objectSpec,
			InputProperties: map[string]schema.PropertySpec{},
			RequiredInputs:  definitionProperties.Required,
		}

		pkgImportAliases[fmt.Sprintf("%s/%s", goImportPath, defintion)] = defintion

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
				if defintion == kindClusterDefinition {
					resourceSpec.Properties["kubeconfig"] = schema.PropertySpec{
						Description: "KubeConfig",
						TypeSpec: schema.TypeSpec{
							Type: "string",
						},
					}
					resourceSpec.Required = []string{"kubeconfig"}
				}
				pkg.Resources[tok] = resourceSpec
				pkg.Types[tok] = schema.ComplexTypeSpec{
					ObjectTypeSpec: schema.ObjectTypeSpec{
						Description: fmt.Sprintf("KIND %s type", definitionPropertyKey),
						Type:        "object",
						Properties:  resourceSpec.InputProperties,
						Language:    map[string]json.RawMessage{},
					},
				}
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

func rawMessage(v interface{}) json.RawMessage {
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
		ref = fmt.Sprintf("%s/kind:index:%s", "#/types", ref)
		return ref
	}
	return ref
}
