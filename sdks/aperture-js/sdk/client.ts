import grpc from "@grpc/grpc-js";
import * as otelApi from "@opentelemetry/api";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-grpc";
import { Resource } from "@opentelemetry/resources";
import { BatchSpanProcessor, Tracer } from "@opentelemetry/sdk-trace-base";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { SemanticResourceAttributes } from "@opentelemetry/semantic-conventions";

import { Flow } from "./flow.js";
import { fcs, FlowControlService } from "./utils.js";
import {
  FLOW_START_TIMESTAMP_LABEL,
  LIBRARY_NAME,
  LIBRARY_VERSION,
  SOURCE_LABEL,
  URL,
  WORKLOAD_START_TIMESTAMP_LABEL,
} from "./consts.js";

/**
 * Types to be exported from package
 */
export type { FlowControlService, Flow };

export { grpc };

type ControlPoint = unknown;
type Label = [string, string];

export class ApertureClient {
  public readonly fcsClient: FlowControlService;

  public readonly exporter: OTLPTraceExporter;

  public readonly tracerProvider: NodeTracerProvider;

  public readonly tracer: Tracer;

  public readonly timeout: number;

  constructor({
    timeout = 200,
    channelCredentials = grpc.credentials.createInsecure(),
  } = {}) {
    this.fcsClient = new fcs.FlowControlService(URL, channelCredentials);

    this.exporter = new OTLPTraceExporter({
      url: URL,
      credentials: channelCredentials,
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

  // StartFlow takes a control point and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
  // Return value is a Flow.
  // The call returns immediately in case connection with Aperture Agent is not established.
  // The default semantics are fail-to-wire. If StartFlow fails, calling Flow.ShouldRun() on returned Flow returns as true.
  async StartFlow(controlPointArg: ControlPoint, labelsArg: Iterable<Label>) {
    return new Promise<Flow>((resolve, reject) => {
      let labelsMap = new Map<string, string>();
      let baggage = otelApi.propagation.getBaggage(otelApi.context.active());

      if (baggage !== undefined) {
        for (const member of baggage.getAllEntries()) {
          labelsMap.set(member[0], member[1].value);
        }
      }

      let mergedLabels = new Map<string, string>([...labelsMap, ...labelsArg]);
      let span = this.tracer.startSpan("Aperture Check");
      span.setAttribute(FLOW_START_TIMESTAMP_LABEL, Date.now());
      span.setAttribute(SOURCE_LABEL, "sdk");
      let flow = new Flow(span);

      let checkParams = { deadline: 0 };
      if (this.timeout != null && this.timeout != 0) {
        checkParams = { deadline: Date.now() + this.timeout };
      }

      this.fcsClient.Check(
        {
          control_point: controlPointArg,
          labels: mergedLabels,
        },
        checkParams,
        (err, response) => {
          span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Date.now());

          if (err) {
            if (flow.failOpen) {
              // Accept the request if failOpen is true even if we encounter an error
              console.log(
                `Aperture server unavailable due to ${JSON.stringify(
                  err
                )}. Accepting request.`
              );
              flow.checkResponse = null;
            } else {
              // Reject the request if failOpen is false if we encounter an error
              reject(err);
            }
          } else {
            flow.checkResponse = response;
          }

          resolve(flow);
        }
      );
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
