package main

import (
	"github.com/frezbo/pulumi-provider-kind/sdk/v3/go/kind/cluster"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cluster, err := cluster.NewCluster(ctx, "kind", &cluster.ClusterArgs{
			Name: pulumi.String("kindest"),
		})
		if err != nil {
			return err
		}
		ctx.Export("name", cluster.Name)
		return nil
	})
}
