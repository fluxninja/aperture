import grpc, { ChannelCredentials, Client } from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";

import { PROTO_PATH } from "./consts.js";
import { Flow } from "./flow.js";
import { Error, Response } from "./types.js";

export type ApertureGrpcObject = {
  aperture: {
    flowcontrol: {
      check: {
        v1: {
          FlowControlService: {
            new (
              url: string,
              channelCredentials: ChannelCredentials,
            ): FlowControlService;
          };
        };
      };
    };
  };
};

type CheckOptions1 = { controlPoint: unknown; labels: unknown };

type CheckOptions2 = {
  deadline: number;
};

type CheckCallback = (err: Error, response: Response) => void;

export abstract class FlowControlService extends Client {
  constructor(url: string, channelCredentials: ChannelCredentials) {
    super(url, channelCredentials);
  }

  public abstract Check(
    options1: CheckOptions1,
    options2: CheckOptions2,
    callback: CheckCallback,
  ): Flow;
}

const clientPackage = protoLoader.loadSync(PROTO_PATH, {
  keepCase: false, // NOTE: make sure we are using camelCase to access the proto
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

export const fcs = (
  grpc.loadPackageDefinition(clientPackage) as unknown as ApertureGrpcObject
).aperture.flowcontrol.check.v1;
