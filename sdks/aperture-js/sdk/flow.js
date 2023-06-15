import {
  CHECK_RESPONSE_LABEL,
  FLOW_END_TIMESTAMP_LABEL,
  FLOW_STATUS_LABEL,
} from "./consts.js";

export const FlowStatus = Object.freeze({
  Ok: "Ok",
  Error: "Error",
});

export const FlowDecision = Object.freeze({
  Accepted: "Accepted",
  Rejected: "Rejected",
  Unreachable: "Unreachable",
})

export class Flow {
  constructor(span, checkResponse = null) {
    this.span = span;
    this.checkResponse = checkResponse;
    this.statusCode = FlowStatus.Ok;
    this.ended = false;
    this.failOpen = true;
  }

  ShouldRun() {
    var decision = this.Decision();
    return decision === FlowDecision.Accepted || (this.failOpen && decision === FlowDecision.Unreachable)
  }

  DisableFailOpen() {
    this.failOpen = false;
  }

  Decision() {
    if (this.checkResponse === undefined) {
      return FlowDecision.Unreachable;
    }
    if (this.checkResponse.decisionType === "DECISION_TYPE_ACCEPTED") {
      return FlowDecision.Accepted;
    }
    return FlowDecision.Rejected;
  }

  SetStatus(statusCode) {
    this.statusCode = statusCode;
  }

  End() {
    if (this.ended) {
      return new Error("flow already ended");
    }
    this.ended = true;

    this.span.setAttribute(FLOW_STATUS_LABEL, this.statusCode);
    this.span.setAttribute(
      CHECK_RESPONSE_LABEL,
      JSON.stringify(this.checkResponse),
    );
    this.span.setAttribute(FLOW_END_TIMESTAMP_LABEL, Date.now());

    let attr = this.span.attributes;
    this.span.end();
  }

  CheckResponse() {
    return this.checkResponse;
  }
}
