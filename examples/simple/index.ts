import * as kind from "../../sdk/nodejs";

const cluster = new kind.Cluster("kind-example", {})

export const kubeconfig = cluster.kubeconfig;
