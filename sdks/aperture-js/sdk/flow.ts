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

export class Flow {
  public ended = false;

  constructor(
    public span: Span,
    public checkResponse: Response | null | undefined = null,
  ) {}

  Accepted() {
    if (this.checkResponse === undefined || this.checkResponse === null) {
      return true;
    }
    if (this.checkResponse.decisionType === "DECISION_TYPE_ACCEPTED") {
      return true;
    }
    return false;
  }

  End(flowStatus: AttributeValue) {
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

    // TODO: attr are unused, can be deleted? MOreover ts throws that attributes do not exist on Span
    // @ts-ignore
    let attr = this.span.attributes;
    this.span.end();
  }

  CheckResponse() {
    return this.checkResponse;
  }
}
