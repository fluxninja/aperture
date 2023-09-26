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

import { LIBRARY_NAME, LIBRARY_VERSION, URL } from "./consts.js";
import { Flow } from "./flow.js";
import { fcs } from "./utils.js";

export interface FlowParams {
  labels?: Record<string, string>;
  timeoutMilliseconds?: number;
  rampMode?: boolean;
  tryConnect?: boolean;
}

export class ApertureClient {
  private readonly fcsClient: FlowControlServiceClient;

  private readonly exporter: OTLPTraceExporter;

  private readonly tracerProvider: NodeTracerProvider;

  private readonly tracer: Tracer;

  constructor({ channelCredentials = grpc.credentials.createInsecure() } = {}) {
    this.fcsClient = new fcs.FlowControlService(URL, channelCredentials, {
      "grpc.keepalive_time_ms": 10000,
      "grpc.keepalive_timeout_ms": 5000,
      "grpc.keepalive_permit_without_calls": 1,
    });

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
  // For FlowParams set defaults - labels = {}, timeoutMilliseconds = 0, rampMode = false using Pick.
  async StartFlow(
    controlPoint: string,
    params: FlowParams = {},
  ): Promise<Flow> {
    return new Promise<Flow>((resolve) => {
      if (params.rampMode === undefined) {
        params.rampMode = false;
      }
      let span = this.tracer.startSpan("Aperture Check");
      let startDate = Date.now();

      const resolveFlow = (response: any, err: any) => {
        resolve(new Flow(span, startDate, params.rampMode, response, err));
      };

      try {
        // check connection state
        // if not ready, return flow with fail-to-wire semantics
        // if ready, call check
        if (
          (params.tryConnect === undefined || params.tryConnect == false) &&
          this.fcsClient.getChannel().getConnectivityState(true) !=
            connectivityState.READY
        ) {
          resolveFlow(null, new Error("connection not ready"));
          return;
        }

        let labelsBaggage = {} as Record<string, string>;
        let baggage = otelApi.propagation.getBaggage(otelApi.context.active());

        if (baggage !== undefined) {
          for (const member of baggage.getAllEntries()) {
            labelsBaggage[member[0]] = member[1].value;
          }
        }

        let mergedLabels = { ...params.labels, ...labelsBaggage };

        const request: CheckRequest = {
          controlPoint: controlPoint,
          labels: mergedLabels,
          rampMode: params.rampMode,
        };

        const grpcParams: grpc.CallOptions = {
          deadline:
            params.timeoutMilliseconds != undefined &&
            params.timeoutMilliseconds > 0
              ? new Date(Date.now() + params.timeoutMilliseconds)
              : undefined,
        };

        const cb: grpc.requestCallback<CheckResponse__Output> = (
          err: any,
          response: any,
        ) => {
          resolveFlow(err ? null : response, err);
          return;
        };

        this.fcsClient.Check(request, grpcParams, cb);
      } catch (err: any) {
        resolveFlow(null, err);
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
