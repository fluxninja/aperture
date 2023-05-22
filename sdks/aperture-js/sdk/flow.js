import {
  CHECK_RESPONSE_LABEL,
  FLOW_END_TIMESTAMP_LABEL,
  FLOW_STATUS_LABEL,
} from "./consts.js";

export const FlowStatus = Object.freeze({
  Ok: "Ok",
  Error: "Error",
});

export const FlowResult = Object.freeze({
  Accepted: "Accepted",
  Rejected: "Rejected",
  Unreachable: "Unreachable",
})

export class Flow {
  constructor(span, checkResponse = null) {
    this.span = span;
    this.ended = false;
    this.checkResponse = checkResponse;
  }

  Result() {
    if (this.checkResponse === undefined) {
      return FlowResult.Unreachable;
    }
    if (this.checkResponse.decisionType === "DECISION_TYPE_ACCEPTED") {
      return FlowResult.Accepted;
    }
    return FlowResult.Rejected;
  }

  End(flowStatus) {
    if (this.ended) {
      return new Error("flow already ended");
    }
    this.ended = true;

    this.span.setAttribute(FLOW_STATUS_LABEL, flowStatus);
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
