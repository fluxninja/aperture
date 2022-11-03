import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import { PROTO_PATH } from "./consts.js";

const clientPackage = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
});

export const fcs = grpc.loadPackageDefinition(clientPackage).aperture.flowcontrol.v1;
