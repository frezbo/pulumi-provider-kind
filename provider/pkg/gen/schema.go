package gen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alecthomas/jsonschema"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/go/common/util/contract"
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

	for defintion, props := range swagger.Definitions {

		defintionLowerCase := strings.ToLower(defintion)
		tok := fmt.Sprintf("kind:%s:%s", defintion, defintion)
		pkgImportAliases[fmt.Sprintf("%s/%s", goImportPath, defintionLowerCase)] = defintionLowerCase

		objectSpec := schema.ObjectTypeSpec{
			Description: fmt.Sprintf("KIND %s", defintion),
			Type:        "object",
			Properties:  map[string]schema.PropertySpec{},
			Language:    map[string]json.RawMessage{},
		}
		pkg.Types[tok] = schema.ComplexTypeSpec{
			ObjectTypeSpec: objectSpec,
		}

		resourceSpec := schema.ResourceSpec{
			ObjectTypeSpec:  objectSpec,
			InputProperties: map[string]schema.PropertySpec{},
			RequiredInputs:  props.Required,
		}

		for _, key := range props.Properties.Keys() {
			if key == "PatchJSON6902" {
				continue
			}
			val, _ := props.Properties.Get(key)
			castedVal := val.(*jsonschema.Type)
			ref := castedVal.Ref
			if ref != "" {
				ref = strings.ReplaceAll(ref, "#/definitions/", "")
				ref = fmt.Sprintf("%s/kind:%s:%s", "#/types", ref, ref)
			}
			inputProps := schema.PropertySpec{
				TypeSpec: schema.TypeSpec{
					Type: castedVal.Type,
					Ref:  ref,
				},
			}
			if castedVal.Items != nil {
				ref := castedVal.Items.Ref
				if ref != "" {
					ref = strings.ReplaceAll(ref, "#/definitions/", "")
					ref = fmt.Sprintf("%s/kind:%s:%s", "#/types", ref, ref)
				}
				inputProps.TypeSpec.Items = &schema.TypeSpec{
					Type: castedVal.Items.Type,
					Ref:  ref,
				}
			}
			resourceSpec.InputProperties[key] = inputProps
			resourceSpec.Properties[key] = inputProps
			pkg.Types[tok].ObjectTypeSpec.Properties[key] = inputProps
		}

		pkg.Resources[tok] = resourceSpec

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
