import { AttributeValue, Span } from "@opentelemetry/api";

import {
  CHECK_RESPONSE_LABEL,
  FLOW_END_TIMESTAMP_LABEL,
  FLOW_STATUS_LABEL,
} from "./consts.js";
import { Response } from "./types.js";

export const FlowStatus = Object.freeze({
  Ok: "Ok",
  Error: "Error",
});

export const FlowDecision = Object.freeze({
  Accepted: "Accepted",
  Rejected: "Rejected",
  Unreachable: "Unreachable",
});

export class Flow {
  constructor(
    public span: Span,
    public checkResponse: Response | null | undefined = null,
    public statusCode: AttributeValue = FlowStatus.Ok,
    public ended: boolean = false,
    public failOpen: boolean = true,
  ) {}

  ShouldRun() {
    var decision = this.Decision();
    return (
      decision === FlowDecision.Accepted ||
      (this.failOpen && decision === FlowDecision.Unreachable)
    );
  }

  DisableFailOpen() {
    this.failOpen = false;
  }

  Decision() {
    if (this.checkResponse === undefined || this.checkResponse === null) {
      return FlowDecision.Unreachable;
    }
    if (this.checkResponse.decisionType === "DECISION_TYPE_ACCEPTED") {
      return FlowDecision.Accepted;
    }
    return FlowDecision.Rejected;
  }

  SetStatus(statusCode: AttributeValue) {
    this.statusCode = statusCode;
  }

  End() {
    if (this.ended) {
      return new Error("flow already ended");
    }
    this.ended = true;

    if (this.checkResponse) {
      // HACK: Change timestamps to ISO strings since the protobufjs library uses it in a different format
      // Issue: https://github.com/protobufjs/protobuf.js/issues/893
      // PR: https://github.com/protobufjs/protobuf.js/pull/1258
      // Current timestamp type: https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/timestamp.proto
      const localCheckResponse = this.checkResponse as any;
      if (
        localCheckResponse.start &&
        typeof localCheckResponse.start === "object"
      ) {
        localCheckResponse.start = new Date(
          localCheckResponse.start.seconds * 1000 +
            localCheckResponse.start.nanos / 1000,
        ).toISOString();
      }
      if (
        localCheckResponse.end &&
        typeof localCheckResponse.end === "object"
      ) {
        localCheckResponse.end = new Date(
          localCheckResponse.end.seconds * 1000 +
            localCheckResponse.end.nanos / 1000,
        ).toISOString();
      }
      this.checkResponse = localCheckResponse;
    }

    this.span.setAttribute(FLOW_STATUS_LABEL, this.statusCode);
    this.span.setAttribute(
      CHECK_RESPONSE_LABEL,
      JSON.stringify(this.checkResponse),
    );
    this.span.setAttribute(FLOW_END_TIMESTAMP_LABEL, Date.now());

    this.span.end();
  }

  CheckResponse() {
    return this.checkResponse;
  }
}
