import grpc from "@grpc/grpc-js";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { Resource } from "@opentelemetry/resources";
import { SemanticResourceAttributes } from "@opentelemetry/semantic-conventions";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-grpc";

import { fcs } from "./utils.js";
import {
  URL,
  LIBRARY_NAME,
  LIBRARY_VERSION,
  FLOW_START_TIMESTAMP_LABEL,
  SOURCE_LABEL,
  WORKLOAD_START_TIMESTAMP_LABEL,
} from "./consts.js";
import { Flow } from "./flow.js";



export class ApertureClient {
  constructor(timeout = 200) {
    this.fcsClient = new fcs.FlowControlService(URL, grpc.credentials.createInsecure());

    let exporter = new OTLPTraceExporter({
      url: URL,
      credentials: grpc.credentials.createInsecure(),
    });
    let res = this.#newResource();
    let tracerProvider = new NodeTracerProvider({
      resource: res,
    });
    tracerProvider.addSpanProcessor(new BatchSpanProcessor(exporter));
    tracerProvider.register();
    this.tracer = tracerProvider.getTracer(LIBRARY_NAME, LIBRARY_VERSION);
  }

  // StartFlow takes a feature name and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
  // Return value is a Flow.
  // The call returns immediately in case connection with Aperture Agent is not established.
  // The default semantics are fail-to-wire. If StartFlow fails, calling Flow.Accepted() on returned Flow returns as true.
  async StartFlow(featureArg, labelsArg) {
    return new Promise((resolve, reject) => {
      // TODO - process baggage

      let span = this.tracer.startSpan("Aperture Check");
      span.setAttributes({
        FLOW_START_TIMESTAMP_LABEL: Date.now(),
        SOURCE_LABEL: "sdk",
      });
      let flow = new Flow(span);

      this.fcsClient.Check(
        {
          feature: featureArg,
          labels: labelsArg,
        }, (err, response) => {
          span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Date.now());

          if (err) {
            if (err.code === grpc.status.UNAVAILABLE) {
              console.log(`Aperture server unavailable. Accepting request.\n`);
              flow.checkResponse = response;
              resolve(flow);
            }
            reject(err);
          }

          console.log(`Response: ${response}\n`);
          flow.checkResponse = response;
          resolve(flow);
        });
    });
  }

  Shutdown() {
    grpc.closeClient(this.fcsClient);
    return;
  }

  GetState() {
    return this.fcsClient.getChannel().getConnectivityState(true);
  }

  #newResource() {
    let defaultRes = Resource.default();
    let res = new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: LIBRARY_NAME,
      [SemanticResourceAttributes.SERVICE_VERSION]: LIBRARY_VERSION,
    });
    let merged = defaultRes.merge(res);
    return merged;
  }
}
