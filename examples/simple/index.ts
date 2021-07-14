import * as kind from "../../sdk/nodejs";

const cluster = new kind.cluster.Cluster("kind-example", {
    nodes: [
        {
            role: kind.node.RoleType.ControlPlane
        },
        {
            role: kind.node.RoleType.Worker
        }
    ],
})

export const kubeconfig = cluster.kubeconfig;
export const name = cluster.name;
