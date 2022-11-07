import grpc from "@grpc/grpc-js";
import * as otelApi from "@opentelemetry/api";
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

    this.exporter = new OTLPTraceExporter({
      url: URL,
      credentials: grpc.credentials.createInsecure(),
    });
    let res = this.#newResource();
    this.tracerProvider = new NodeTracerProvider({
      resource: res,
    });
    this.tracerProvider.addSpanProcessor(new BatchSpanProcessor(this.exporter));
    this.tracerProvider.register();
    this.tracer = this.tracerProvider.getTracer(LIBRARY_NAME, LIBRARY_VERSION);
    this.timeout = timeout;
  }

  // StartFlow takes a feature name and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
  // Return value is a Flow.
  // The call returns immediately in case connection with Aperture Agent is not established.
  // The default semantics are fail-to-wire. If StartFlow fails, calling Flow.Accepted() on returned Flow returns as true.
  async StartFlow(featureArg, labelsArg) {
    return new Promise((resolve, reject) => {
      let labelsMap = new Map();
      let baggage = otelApi.propagation.getBaggage(otelApi.context.active());
      if (baggage !== undefined) {
        for (const member of baggage.getAllEntries()) {
          labelsMap[member[0]] = member[1].value;
        }
      }

      let mergedLabels = new Map([...labelsMap, ...labelsArg])
      let span = this.tracer.startSpan("Aperture Check");
      span.setAttribute(FLOW_START_TIMESTAMP_LABEL, Date.now());
      span.setAttribute(SOURCE_LABEL, "sdk");
      let flow = new Flow(span);

      this.fcsClient.Check(
        {
          feature: featureArg,
          labels: mergedLabels,
        },
        { deadline: Date.now() + this.timeout },
        (err, response) => {
          span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Date.now());

          if (err) {
            if (err.code === grpc.status.UNAVAILABLE) {
              console.log(`Aperture server unavailable. Accepting request.`);
              resolve(flow);
            }
            reject(err);
          }

          flow.checkResponse = response;
          resolve(flow);
        });
    });
  }

  Shutdown() {
    grpc.closeClient(this.fcsClient);
    this.exporter.shutdown();
    this.tracerProvider.shutdown();
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
