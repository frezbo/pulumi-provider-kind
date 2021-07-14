package main

import (
	"github.com/frezbo/pulumi-provider-kind/sdk/v3/go/kind/cluster"
	"github.com/frezbo/pulumi-provider-kind/sdk/v3/go/kind/node"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// create a simple two node kind cluster
		// with one control-plane and one worer node
		cluster, err := cluster.NewCluster(ctx, "kind-example", &cluster.ClusterArgs{
			Nodes: node.NodeArray{
				node.NodeArgs{
					Role: node.RoleTypeControlPlane,
				},
				node.NodeArgs{
					Role: node.RoleTypeWorker,
				},
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("kubeconfig", cluster.Kubeconfig)
		return nil
	})
}
