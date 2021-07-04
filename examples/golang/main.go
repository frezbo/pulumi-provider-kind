package main

import (
	"github.com/frezbo/pulumi-provider-kind/sdk/v3/go/kind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cluster, err := kind.NewCluster(ctx, "kind-example", &kind.ClusterArgs{
			Networking: kind.NetworkingTypeArgs{
				DisableDefaultCNI: pulumi.Bool(true),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("kubeconfig", cluster.Kubeconfig)
		return nil
	})
}
