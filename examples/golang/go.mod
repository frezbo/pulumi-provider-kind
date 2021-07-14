module kind-go-example

go 1.16

require (
	github.com/frezbo/pulumi-provider-kind/sdk/v3 v3.0.0-20210705082738-d13d4607d669
	github.com/pulumi/pulumi/sdk/v3 v3.7.0
)

replace github.com/frezbo/pulumi-provider-kind/sdk/v3 => ../../sdk
