import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";

import { PROTO_PATH, URL } from "./consts.js";

const clientPackage = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const fcs = grpc.loadPackageDefinition(clientPackage).aperture.flowcontrol.v1;

export class ApertureClient {
  me = new fcs.FlowControlService(URL, grpc.credentials.createInsecure());

  async StartFlow(featureArg, labelsArg) {
    return new Promise((resolve, reject) => {
      this.me.Check(
      {
        feature: featureArg,
        labels: labelsArg,
      }, (err, response) => {
        if (err) {
          reject(err);
        }
        console.log("Response ", response);
        resolve(response);
      });
    });
  }

  async Shutdown() {
    return new Promise(() => {
      // TODO
    });
  }

  async GetState() {
    return new Promise((resolve, reject) => {
      if ( this.me.getChannel().getConnectivityState(true) ) {
        console.log("connected");
        resolve(true);
      }
      else {
        reject();
      }
    });
  }
}
