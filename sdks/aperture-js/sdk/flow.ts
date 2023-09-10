import { Span } from "@opentelemetry/api";

import {
  CHECK_RESPONSE_LABEL,
  FLOW_END_TIMESTAMP_LABEL,
  FLOW_START_TIMESTAMP_LABEL,
  FLOW_STATUS_LABEL,
  SOURCE_LABEL,
  WORKLOAD_START_TIMESTAMP_LABEL,
} from "./consts.js";
import {
  CheckResponse__Output,
  _aperture_flowcontrol_check_v1_CheckResponse_DecisionType,
} from "./gen/aperture/flowcontrol/check/v1/CheckResponse.js";

export const FlowStatusEnum = {
  OK: "OK",
  Error: "Error",
} as const;

export type FlowStatus = (typeof FlowStatusEnum)[keyof typeof FlowStatusEnum];

export class Flow {
  private ended: boolean = false;
  private status: FlowStatus = FlowStatusEnum.OK;

  constructor(
    private span: Span,
    startDate: number,
    private failOpen: boolean = true,
    private checkResponse: CheckResponse__Output | null = null,
    private error: Error | null = null,
  ) {
    span.setAttribute(SOURCE_LABEL, "sdk");
    span.setAttribute(FLOW_START_TIMESTAMP_LABEL, startDate * 1000);
    span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Date.now() * 1000);
  }

  ShouldRun() {
    if (
      (this.failOpen && this.checkResponse === null) ||
      this.checkResponse?.decisionType ===
        _aperture_flowcontrol_check_v1_CheckResponse_DecisionType.DECISION_TYPE_ACCEPTED
    ) {
      return true;
    } else {
      return false;
    }
  }

  SetStatus(status: FlowStatus) {
    this.status = status;
  }

  Error() {
    return this.error;
  }

  CheckResponse() {
    return this.checkResponse;
  }

  Span() {
    return this.span;
  }

  End() {
    if (this.ended) {
      return;
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

      this.span.setAttribute(
        CHECK_RESPONSE_LABEL,
        JSON.stringify(localCheckResponse),
      );
    }

    this.span.setAttribute(FLOW_STATUS_LABEL, this.status);

    this.span.setAttribute(FLOW_END_TIMESTAMP_LABEL, Date.now() * 1000);

    this.span.end();
  }
}
