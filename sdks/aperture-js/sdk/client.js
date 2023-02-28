import grpc from "@grpc/grpc-js";
import * as otelApi from "@opentelemetry/api";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-grpc";
import { Resource } from "@opentelemetry/resources";
import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { SemanticResourceAttributes } from "@opentelemetry/semantic-conventions";

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

/**
 * Class representing Aperture client.
 * @class
 */
export class ApertureClient {
  /**
   * Create an Aperture client.
   * @constructor
   * @param {number} [timeout=200] - Timeout for the client in milliseconds.
   */
  constructor(timeout = 200) {
    this.fcsClient = new fcs.FlowControlService(
      URL,
      grpc.credentials.createInsecure(),
    );

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

  /**
   * Starts a new flow with given control point and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
   * The call returns immediately in case connection with Aperture Agent is not established.
   * The default semantics are fail-to-wire. If StartFlow fails, calling Flow.Accepted() on returned Flow returns as true.
   * @param {string} controlPointArg - A control point.
   * @param {LabelMap[]} labelsArg - Labels for the flow.
   * @returns {Promise<Flow>} - A new Flow.
   */
  async StartFlow(controlPointArg, labelsArg) {
    return new Promise((resolve, reject) => {
      let labelsMap = new Map();
      let baggage = otelApi.propagation.getBaggage(otelApi.context.active());
      if (baggage !== undefined) {
        for (const member of baggage.getAllEntries()) {
          labelsMap[member[0]] = member[1].value;
        }
      }

      let mergedLabels = new Map([...labelsMap, ...labelsArg]);
      let span = this.tracer.startSpan("Aperture Check");
      span.setAttribute(FLOW_START_TIMESTAMP_LABEL, Date.now());
      span.setAttribute(SOURCE_LABEL, "sdk");
      let flow = new Flow(span);

      this.fcsClient.Check(
        {
          control_point: controlPointArg,
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
        },
      );
    });
  }

  Shutdown() {
    grpc.closeClient(this.fcsClient);
    this.exporter.shutdown();
    this.tracerProvider.shutdown();
    return;
  }

  /**
   * Returns the connectivity state of the FcsClient channel.
   * @returns {grpc.connectivityState}
   */
  GetState() {
    return this.fcsClient.getChannel().getConnectivityState(true);
  }

  /**
   * Creates a new resource object with default and library-specific attributes.
   * @private
   * @returns {Resource}
   */
  #newResource() {
    let defaultRes = Resource.default();
    let res = new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: LIBRARY_NAME,
      [SemanticResourceAttributes.SERVICE_VERSION]: LIBRARY_VERSION,
    });
    let merged = defaultRes.merge(res);
    return merged;
  }

  /**
   * @callback middlewareCallback
   * @param {Object} req - Request object.
   * @param {Object} res - Response object.
   * @param {Function} next - Next middleware function.
   * @returns {void}
   */

  /**
   * Returns middleware for express.js.
   * @param {string} [controlPoint] - A control point - defaults to req.path.
   * @param {LabelMap[]} [labels] - Labels for the flow.
   * @param {number} [timeout] - Timeout for the client in milliseconds.
   * @return {middlewareCallback}
   */
  Middleware(controlPoint, labels, timeout) {
    if (timeout === undefined) {
      timeout = this.timeout;
    }
    if (labels === undefined) {
      labels = new Map();
    }
    return (req, res, next) => {
      if (controlPoint === undefined) {
        controlPoint = req.path;
      }
      labels = new Map(labels);
      labels.set("http_method", req.method);
      labels.set("http_host", req.headers.host);
      labels.set("http_path", req.path);
      labels.set("http_user_agent", req.headers["user-agent"]);
      labels.set("http_content_length", req.headers["content-length"]);
      this.StartFlow(controlPoint, labels)
        .then((flow) => {
          if (flow.Accepted()) {
            next();
            flow.End(FlowStatus.Ok);
          } else {
            res.status(403).send("Forbidden");
            flow.End(FlowStatus.Error);
          }
        })
        .catch(next);
    };
  }
}
