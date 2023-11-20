import grpc from "@grpc/grpc-js";
import path from "path";
import protoLoader from "@grpc/proto-loader";
import { ProtoGrpcType } from "./gen/check";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);

export const PROTO_PATH = path.resolve(
  path.dirname(__filename),
  "../proto/flowcontrol/check/v1/check.proto",
);

const clientPackage = protoLoader.loadSync(PROTO_PATH, {
  defaults: true,
  longs: String,
});

export const checkv1 = (grpc.loadPackageDefinition(clientPackage) as unknown as ProtoGrpcType).aperture.flowcontrol.check.v1;
