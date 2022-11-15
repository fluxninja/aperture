import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import { PROTO_PATH } from "./consts.js";

const clientPackage = protoLoader.loadSync(PROTO_PATH, {
    keepCase: false,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
});

export const fcs = grpc.loadPackageDefinition(clientPackage).aperture.flowcontrol.check.v1;
