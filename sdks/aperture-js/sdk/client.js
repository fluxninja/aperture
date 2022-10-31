import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";

import { PROTO_PATH, URL } from "./consts.js";
import { Flow } from "./flow.js";

const clientPackage = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const fcs = grpc.loadPackageDefinition(clientPackage).aperture.flowcontrol.v1;

export class ApertureClient {
  constructor() {
    this.fcsClient = new fcs.FlowControlService(URL, grpc.credentials.createInsecure());
  }

  async StartFlow(featureArg, labelsArg) {
    return new Promise((resolve, reject) => {
      this.fcsClient.Check(
        {
          feature: featureArg,
          labels: labelsArg,
        }, (err, response) => {
          if (err) {
            reject(err);
          }

          console.log(`Response: ${response}\n`);
          let flow = new Flow(response);
          resolve(flow);
        });
    });
  }

  Shutdown() {
    // TODO
    return;
  }

  GetState() {
    return this.fcsClient.getChannel().getConnectivityState(true);
  }
}
