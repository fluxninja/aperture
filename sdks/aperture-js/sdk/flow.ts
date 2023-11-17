import { Span } from "@opentelemetry/api";
import {
  CHECK_RESPONSE_LABEL,
  FLOW_END_TIMESTAMP_LABEL,
  FLOW_START_TIMESTAMP_LABEL,
  FLOW_STATUS_LABEL,
  SOURCE_LABEL,
  WORKLOAD_START_TIMESTAMP_LABEL,
} from "./consts.js";
import grpc from "@grpc/grpc-js";
import {
  CheckResponse__Output,
  _aperture_flowcontrol_check_v1_CheckResponse_DecisionType,
} from "./gen/aperture/flowcontrol/check/v1/CheckResponse.js";
import type { Duration__Output as _google_protobuf_Duration__Output } from "./gen/google/protobuf/Duration";
import type { Timestamp__Output as _google_protobuf_Timestamp__Output } from "./gen/google/protobuf/Timestamp";
import { FlowControlServiceClient } from "./gen/aperture/flowcontrol/check/v1/FlowControlService.js";
import type { CacheUpsertRequest } from "./gen/aperture/flowcontrol/check/v1/CacheUpsertRequest";
import type { CacheDeleteRequest } from "./gen/aperture/flowcontrol/check/v1/CacheDeleteRequest.js";
import { Duration } from "@grpc/grpc-js/build/src/duration.js";
import { CachedValueResponse, SetCachedValueResponse, DeleteCachedValueResponse, } from "./cache.js";

export const FlowStatusEnum = {
  OK: "OK",
  Error: "Error",
} as const;

export type FlowStatus = (typeof FlowStatusEnum)[keyof typeof FlowStatusEnum];

export class Flow {
  private ended: boolean = false;
  private status: FlowStatus = FlowStatusEnum.OK;

  constructor(
    private fcsClient: FlowControlServiceClient,
    private grpcCallOptions: grpc.CallOptions,
    private controlPoint: string,
    private span: Span,
    startDate: number,
    private rampMode: boolean = false,
    private cacheKey: string | null = null,
    private checkResponse: CheckResponse__Output | null = null,
    private error: Error | null = null,
  ) {
    span.setAttribute(SOURCE_LABEL, "sdk");
    span.setAttribute(FLOW_START_TIMESTAMP_LABEL, startDate);
    span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Date.now());
  }

  ShouldRun() {
    if (
      (!this.rampMode && this.checkResponse === null) ||
      (this.checkResponse?.decisionType === _aperture_flowcontrol_check_v1_CheckResponse_DecisionType.DECISION_TYPE_ACCEPTED)
    ) {
      return true;
    } else {
      return false;
    }
  }

  SetStatus(status: FlowStatus) {
    this.status = status;
  }

  async SetCachedValue(value: Buffer, ttl: Duration) {
    if (!this.cacheKey) {
      return Promise.reject(new Error("No cache key"));
    }

    const key = this.cacheKey;
    return new Promise<SetCachedValueResponse | undefined>((resolve) => {
      const cacheUpsertRequest: CacheUpsertRequest = {
        controlPoint: this.controlPoint,
        key: key,
        value: value,
        ttl: ttl,
      }
      this.fcsClient.CacheUpsert(cacheUpsertRequest, this.grpcCallOptions, (err, res) => {
        const resp: SetCachedValueResponse = {
          error: err ?? null,
          code: res?.code.toString() ?? null,
          message: res?.message ?? null,
        };
        resolve(resp);
      });
    });
  }

  async DeleteCachedValue() {
    if (!this.cacheKey) {
      return Promise.reject(new Error("No cache key"));
    }

    const key = this.cacheKey;
    return new Promise<DeleteCachedValueResponse | undefined>((resolve, reject) => {
      const cacheDeleteRequest: CacheDeleteRequest = {
        controlPoint: this.controlPoint,
        key: key,
      }
      this.fcsClient.CacheDelete(cacheDeleteRequest, this.grpcCallOptions, (err, res) => {
        const resp: DeleteCachedValueResponse = {
          error: err ?? null,
          code: res?.code.toString() ?? null,
          message: res?.message ?? null,
        };
        resolve(resp);
      });
    });
  }

  CachedValue() {
    const resp: CachedValueResponse = {
      error: this.error ?? null,
      lookupResult: this.checkResponse?.cachedValue?.lookupResult.toString() ?? null,
      code: this.checkResponse?.cachedValue?.responseCode.toString() ?? null,
      message: this.checkResponse?.cachedValue?.message ?? null,
      value: this.checkResponse?.cachedValue?.value ?? null,
    };
    return resp;
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

      localCheckResponse.start = this.protoTimestampToJSON(this.checkResponse.start);
      localCheckResponse.end = this.protoTimestampToJSON(this.checkResponse.end);
      localCheckResponse.waitTime = this.protoDurationToJSON(this.checkResponse.waitTime);

      // Walk through individual decisions and convert waitTime fields,
      // then add to localCheckResponse, preserving immutability.
      if (this.checkResponse.limiterDecisions) {
        const decisions = this.checkResponse.limiterDecisions.map(
          (decision) => {
            return {
              ...decision,
              waitTime: this.protoDurationToJSON(decision.waitTime),
            };
          },
        );
        localCheckResponse.limiterDecisions = decisions;
      }

      this.span.setAttribute(CHECK_RESPONSE_LABEL, JSON.stringify(localCheckResponse));
    }

    this.span.setAttribute(FLOW_STATUS_LABEL, this.status);
    this.span.setAttribute(FLOW_END_TIMESTAMP_LABEL, Date.now());
    this.span.end();
  }

  private protoTimestampToJSON(timestamp: _google_protobuf_Timestamp__Output | null) {
    if (timestamp) {
      return new Date(
        Number(timestamp.seconds) * 1000 + timestamp.nanos / 1000000,
      ).toISOString();
    }
    return timestamp;
  }

  private protoDurationToJSON(duration: _google_protobuf_Duration__Output | null) {
    if (duration) {
      return `${duration.seconds}.${duration.nanos}s`;
    }
    return duration;
  }
}
