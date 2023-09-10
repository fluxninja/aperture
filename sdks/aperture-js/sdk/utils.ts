import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";

import { PROTO_PATH } from "./consts.js";

import { ProtoGrpcType } from "./gen/check.js";

const clientPackage = protoLoader.loadSync(PROTO_PATH, {
  defaults: true,
  longs: String,
});

export const fcs = (
  grpc.loadPackageDefinition(clientPackage) as unknown as ProtoGrpcType
).aperture.flowcontrol.check.v1;
