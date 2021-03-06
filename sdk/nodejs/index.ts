// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "./utilities";

// Export members:
export * from "./provider";

// Export sub-modules:
import * as cluster from "./cluster";
import * as config from "./config";
import * as node from "./node";
import * as types from "./types";

export {
    cluster,
    config,
    node,
    types,
};

import { Provider } from "./provider";

pulumi.runtime.registerResourcePackage("kind", {
    version: utilities.getVersion(),
    constructProvider: (name: string, type: string, urn: string): pulumi.ProviderResource => {
        if (type !== "pulumi:providers:kind") {
            throw new Error(`unknown provider type ${type}`);
        }
        return new Provider(name, <any>undefined, { urn });
    },
});
