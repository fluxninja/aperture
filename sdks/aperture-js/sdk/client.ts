import grpc, {
  ChannelCredentials,
  ChannelOptions,
  connectivityState,
} from "@grpc/grpc-js";
import * as otelApi from "@opentelemetry/api";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-grpc";
import { Resource } from "@opentelemetry/resources";
import { BatchSpanProcessor, Tracer } from "@opentelemetry/sdk-trace-base";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { SemanticResourceAttributes } from "@opentelemetry/semantic-conventions";
import { LIBRARY_NAME, LIBRARY_VERSION } from "./consts.js";
import { Flow, _Flow } from "./flow.js";
import { CheckRequest } from "./gen/aperture/flowcontrol/check/v1/CheckRequest.js";
import { CheckResponse__Output } from "./gen/aperture/flowcontrol/check/v1/CheckResponse.js";
import { FlowControlServiceClient } from "./gen/aperture/flowcontrol/check/v1/FlowControlService.js";
import { fcs } from "./utils.js";

/**
 * Represents the parameters for a flow.
 */
export interface FlowParams {
  /**
   * Optional labels for the flow.
   */
  labels?: Record<string, string>;
  /**
   * Specifies whether the flow should use ramp mode.
   */
  rampMode?: boolean;
  /**
   * Additional gRPC call options for the flow.
   */
  grpcCallOptions?: grpc.CallOptions;
  /**
   * Specifies whether to try connecting to the flow.
   */
  tryConnect?: boolean;
  /**
   * Key to the result cache entry which needs to be fetched at flow start.
   */
  resultCacheKey?: string;
  /**
   * Keys to global cache entries that need to be fetched at flow start.
   */
  globalCacheKeys?: string[];
}

/**
 * Represents the Aperture Client used for interacting with the Aperture Agent.
 * @example
 * ```ts
 *const apertureClient = new ApertureClient({
 *  address:
 *    process.env.APERTURE_AGENT_ADDRESS !== undefined
 *      ? process.env.APERTURE_AGENT_ADDRESS
 *      : "localhost:8089",
 *  apiKey: process.env.APERTURE_API_KEY || undefined,
 *  // if process.env.APERTURE_AGENT_INSECURE set channelCredentials to insecure
 *  channelCredentials:
 *    process.env.APERTURE_AGENT_INSECURE !== undefined
 *      ? grpc.credentials.createInsecure()
 *      : grpc.credentials.createSsl(),
 *});
 * ```
 */
export class ApertureClient {
  private readonly fcsClient: FlowControlServiceClient;

  private readonly exporter: OTLPTraceExporter;

  private readonly tracerProvider: NodeTracerProvider;

  private readonly tracer: Tracer;

  /**
   * Constructs a new instance of the ApertureClient.
   * @param address The address of the Aperture Agent.
   * @param agentAPIKey The API key for the Aperture Agent (optional).
   * @param channelCredentials The credentials for the gRPC channel (optional).
   * @param channelOptions The options for the gRPC channel (optional).
   * @throws Error if the address is not provided.
   */
  constructor({
    address,
    apiKey,
    channelCredentials = grpc.credentials.createSsl(),
    channelOptions = {},
  }: {
    address: string;
    channelCredentials?: ChannelCredentials;
    channelOptions?: ChannelOptions;
    apiKey?: string;
  }) {
    if (!address) {
      throw new Error("address is required");
    }

    if (apiKey) {
      channelCredentials = grpc.credentials.combineChannelCredentials(
        channelCredentials,
        grpc.credentials.createFromMetadataGenerator(
          (_params: any, callback: any) => {
            const metadata = new grpc.Metadata();
            metadata.add("x-api-key", apiKey);
            callback(null, metadata);
          },
        ),
      );
    }

    this.fcsClient = new fcs.FlowControlService(
      address,
      channelCredentials,
      channelOptions,
    );

    this.exporter = new OTLPTraceExporter({
      url: address,
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

  /**
   * Starts a new flow with the specified control point and parameters.
   * startFlow() takes a control point and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
   * Return value is a Flow.
   * The default semantics are fail-to-wire. If startFlow() fails, calling Flow.ShouldRun() on returned Flow returns as true.
   * @param controlPoint The control point for the flow.
   * @param params The parameters for the flow.
   * @returns A promise that resolves to a Flow object.
   * @example
   * ```ts
   *apertureClient.startFlow("awesomeFeature", {
   *  labels: labels,
   *  grpcCallOptions: {
   *    deadline: Date.now() + 30000,
   *  },
   *  rampMode: false,
   *  cacheKey: "cache",
   *});
   */
  async startFlow(controlPoint: string, params: FlowParams): Promise<Flow> {
    return new Promise<Flow>((resolve) => {
      if (params.rampMode === undefined) {
        params.rampMode = false;
      }
      let span = this.tracer.startSpan("Aperture Check");
      let startDate = Date.now();

      const resolveFlow = (response: any, err: any) => {
        resolve(
          new _Flow(
            this.fcsClient,
            controlPoint,
            span,
            startDate,
            params.rampMode,
            params.resultCacheKey,
            response,
            err,
          ),
        );
      };

      try {
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
          cacheLookupRequest: {
            resultCacheKey: params.resultCacheKey,
            globalCacheKeys: params.globalCacheKeys,
          },
          expectEnd: true,
        };

        const cb: grpc.requestCallback<CheckResponse__Output> = (
          err: any,
          response: any,
        ) => {
          resolveFlow(err ? null : response, err);
          return;
        };

        if (params.grpcCallOptions === undefined) {
          params.grpcCallOptions = {};
        }

        this.fcsClient.Check(request, params.grpcCallOptions, cb);
      } catch (err: any) {
        resolveFlow(null, err);
      }
    });
  }

  /**
   * Shuts down the ApertureClient.
   */
  shutdown() {
    this.fcsClient.getChannel().close();
    this.exporter.shutdown();
    this.tracerProvider.shutdown();
    return;
  }

  /**
   * Gets the current state of the gRPC channel.
   * @returns The connectivity state of the channel.
   */
  getState() {
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
