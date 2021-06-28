module kind-go-example

go 1.16

require (
	github.com/frezbo/pulumi-provider-kind/sdk/v3 v3.0.0-20210613163246-118874c291e9
	github.com/pulumi/pulumi/sdk/v3 v3.5.1
)

replace github.com/frezbo/pulumi-provider-kind/sdk/v3 => ../../sdk
