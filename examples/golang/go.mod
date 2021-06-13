module kind-go-example

go 1.16

require (
	github.com/frezbo/pulumi-provider-kind/sdk/v3 v3.0.0-00010101000000-000000000000
	github.com/pulumi/pulumi/sdk/v3 v3.4.0
)

replace github.com/frezbo/pulumi-provider-kind/sdk/v3 => ../../sdk
