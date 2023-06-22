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
})

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
    return decision === FlowDecision.Accepted || (this.failOpen && decision === FlowDecision.Unreachable)
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

    this.span.setAttribute(FLOW_STATUS_LABEL, this.statusCode);
    this.span.setAttribute(
      CHECK_RESPONSE_LABEL,
      JSON.stringify(this.checkResponse),
    );
    this.span.setAttribute(FLOW_END_TIMESTAMP_LABEL, Date.now());

    // TODO: attr are unused, can be deleted? MOreover ts throws that attributes do not exist on Span
    // @ts-ignore
    let attr = this.span.attributes;
    this.span.end();
  }

  CheckResponse() {
    return this.checkResponse;
  }
}
