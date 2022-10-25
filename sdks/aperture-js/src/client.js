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

class Client {
  me = new fcs.FlowControlService(URL, grpc.credentials.createInsecure());

  async StartFlow(featureArg, labelsArg) {
    var checkRequest = fcs.CheckRequest({
      feature: featureArg,
      labels: labelsArg,
    });

    return new Promise((resolve, reject) => {
      this.me.Check(checkRequest, (err, response) => {
        if (err) {
          console.log("Got error ", err);
          reject(err);
        }
        console.log("Response ", response);
        resolve(response);
      });
    });
  }
}

export const client = new Client();
