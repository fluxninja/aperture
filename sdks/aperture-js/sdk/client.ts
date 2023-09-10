import grpc, { connectivityState } from "@grpc/grpc-js";
import * as otelApi from "@opentelemetry/api";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-grpc";
import { Resource } from "@opentelemetry/resources";
import { BatchSpanProcessor, Tracer } from "@opentelemetry/sdk-trace-base";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { SemanticResourceAttributes } from "@opentelemetry/semantic-conventions";
import { CheckRequest } from "./gen/aperture/flowcontrol/check/v1/CheckRequest.js";
import { CheckResponse__Output } from "./gen/aperture/flowcontrol/check/v1/CheckResponse.js";
import { FlowControlServiceClient } from "./gen/aperture/flowcontrol/check/v1/FlowControlService.js";

import {
  FLOW_START_TIMESTAMP_LABEL,
  LIBRARY_NAME,
  LIBRARY_VERSION,
  SOURCE_LABEL,
  URL,
  WORKLOAD_START_TIMESTAMP_LABEL,
} from "./consts.js";
import { Flow } from "./flow.js";
import { fcs } from "./utils.js";

export class ApertureClient {
  private readonly fcsClient: FlowControlServiceClient;

  private readonly exporter: OTLPTraceExporter;

  private readonly tracerProvider: NodeTracerProvider;

  private readonly tracer: Tracer;

  // Timeout is duration in milliseconds.
  private readonly timeoutMilliseconds: number;

  constructor({
    timeoutMilliseconds = 200,
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

    this.timeoutMilliseconds = timeoutMilliseconds;

    const kickChannel = () => {
      const state = this.fcsClient.getChannel().getConnectivityState(true);
      if (state != connectivityState.SHUTDOWN) {
        this.fcsClient
          .getChannel()
          .watchConnectivityState(state, Infinity, kickChannel);
      }
    };
    kickChannel();
  }

  // StartFlow takes a control point and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
  // Return value is a Flow.
  // The call returns immediately in case connection with Aperture Agent is not established.
  // The default semantics are fail-to-wire. If StartFlow fails, calling Flow.ShouldRun() on returned Flow returns as true.
  async StartFlow(
    controlPointArg: string,
    labels: Record<string, string> = {},
    failOpen: boolean = true,
  ): Promise<Flow> {
    return new Promise<Flow>((resolve) => {
      let span = this.tracer.startSpan("Aperture Check");
      span.setAttribute(FLOW_START_TIMESTAMP_LABEL, Date.now() * 1000);
      span.setAttribute(SOURCE_LABEL, "sdk");

      try {
        // check connection state
        // if not ready, return flow with fail-to-wire semantics
        // if ready, call check
        if (
          this.fcsClient.getChannel().getConnectivityState(true) !=
          connectivityState.READY
        ) {
          resolve(
            new Flow(span, failOpen, null, new Error("connection not ready")),
          );
          return;
        }

        let labelsBaggage = {} as Record<string, string>;
        let baggage = otelApi.propagation.getBaggage(otelApi.context.active());

        if (baggage !== undefined) {
          for (const member of baggage.getAllEntries()) {
            labelsBaggage[member[0]] = member[1].value;
          }
        }

        let mergedLabels = { ...labels, ...labelsBaggage };

        let checkParams: grpc.CallOptions = {};
        if (this.timeoutMilliseconds != 0) {
          const deadline = new Date();
          deadline.setMilliseconds(
            deadline.getMilliseconds() + this.timeoutMilliseconds,
          );
          checkParams.deadline = deadline;
        }

        this.fcsClient.Check(
          {
            controlPoint: controlPointArg,
            labels: mergedLabels,
          } as CheckRequest,
          checkParams as grpc.CallOptions,
          ((err, response) => {
            resolve(new Flow(span, failOpen, response, err));
            return;
          }) as grpc.requestCallback<CheckResponse__Output>,
        );
      } catch (err: any) {
        resolve(new Flow(span, failOpen, null, err));
        return;
      } finally {
        span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Date.now() * 1000);
      }
    });
  }

  Shutdown() {
    this.fcsClient.getChannel().close();
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
