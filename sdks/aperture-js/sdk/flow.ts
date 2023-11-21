import grpc from "@grpc/grpc-js";
import { Duration } from "@grpc/grpc-js/build/src/duration";
import { Span } from "@opentelemetry/api";
import {
  CachedValueResponse,
  ConvertCacheError,
  ConvertCacheLookupStatus,
  ConvertCacheOperationStatus,
  DeleteCachedValueResponse,
  LookupStatus,
  OperationStatus,
  SetCachedValueResponse,
} from "./cache";
import {
  CHECK_RESPONSE_LABEL,
  FLOW_END_TIMESTAMP_LABEL,
  FLOW_START_TIMESTAMP_LABEL,
  FLOW_STATUS_LABEL,
  SOURCE_LABEL,
  WORKLOAD_START_TIMESTAMP_LABEL,
} from "./consts";
import type { CacheDeleteRequest } from "./gen/aperture/flowcontrol/check/v1/CacheDeleteRequest";
import type { CacheUpsertRequest } from "./gen/aperture/flowcontrol/check/v1/CacheUpsertRequest";
import {
  CheckResponse__Output,
  _aperture_flowcontrol_check_v1_CheckResponse_DecisionType,
} from "./gen/aperture/flowcontrol/check/v1/CheckResponse";
import { FlowControlServiceClient } from "./gen/aperture/flowcontrol/check/v1/FlowControlService";
import type { Duration__Output as _google_protobuf_Duration__Output } from "./gen/google/protobuf/Duration";
import type { Timestamp__Output as _google_protobuf_Timestamp__Output } from "./gen/google/protobuf/Timestamp";


/**
 * Enum representing the status of a flow.
 */
export const FlowStatusEnum = {
  OK: "OK",
  Error: "Error",
} as const;

/**
 * Represents the status of a flow.
 */
export type FlowStatus = (typeof FlowStatusEnum)[keyof typeof FlowStatusEnum];

/**
 * Represents a flow in the SDK.
 */
export class Flow {
  private ended: boolean = false;
  private status: FlowStatus = FlowStatusEnum.OK;

  constructor(
    private fcsClient: FlowControlServiceClient,
    private grpcCallOptions: grpc.CallOptions,
    private controlPoint: string,
    private span: Span,
    private startDate: number,
    private rampMode: boolean = false,
    private cacheKey: string | null = null,
    private checkResponse: CheckResponse__Output | null = null,
    private error: Error | null = null,
  ) {
    span.setAttribute(SOURCE_LABEL, "sdk");
    span.setAttribute(FLOW_START_TIMESTAMP_LABEL, this.startDate);
    span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Date.now());
  }

  /**
   * Determines whether the flow should run based on the check response and ramp mode.
   * @returns A boolean value indicating whether the flow should run.
   */
  ShouldRun() {
    if (
      (!this.rampMode && this.checkResponse === null) ||
      this.checkResponse?.decisionType ===
      _aperture_flowcontrol_check_v1_CheckResponse_DecisionType.DECISION_TYPE_ACCEPTED
    ) {
      return true;
    } else {
      return false;
    }
  }

  /**
   * Sets the status of the flow.
   * @param status The status to set.
   */
  SetStatus(status: FlowStatus) {
    this.status = status;
  }

  /**
   * Sets the cached value for the flow.
   * @param value The value to set.
   * @param ttl The time-to-live for the cached value.
   * @returns A promise that resolves to the response of the cache upsert operation.
   */
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
      };
      this.fcsClient.CacheUpsert(
        cacheUpsertRequest,
        this.grpcCallOptions,
        (err, res) => {
          if (err) {
            const resp = new SetCachedValueResponse(err, OperationStatus.Error);
            resolve(resp);
            return;
          }
          const resp = new SetCachedValueResponse(
            ConvertCacheError(res?.error),
            ConvertCacheOperationStatus(res?.operationStatus),
          );
          resolve(resp);
        },
      );
    });
  }

  /**
   * Deletes the cached value for the flow.
   * @returns A promise that resolves to the response of the cache delete operation.
   */
  async DeleteCachedValue() {
    if (!this.cacheKey) {
      return Promise.reject(new Error("No cache key"));
    }

    const key = this.cacheKey;
    return new Promise<DeleteCachedValueResponse | undefined>(
      (resolve, reject) => {
        const cacheDeleteRequest: CacheDeleteRequest = {
          controlPoint: this.controlPoint,
          key: key,
        };
        this.fcsClient.CacheDelete(
          cacheDeleteRequest,
          this.grpcCallOptions,
          (err, res) => {
            if (err) {
              const resp = new DeleteCachedValueResponse(
                err,
                OperationStatus.Error,
              );
              resolve(resp);
              return;
            }
            const resp = new DeleteCachedValueResponse(
              ConvertCacheError(res?.error),
              ConvertCacheOperationStatus(res?.operationStatus),
            );
            resolve(resp);
          },
        );
      },
    );
  }

  /**
   * Gets the cached value for the flow.
   * @returns The cached value response.
   */
  CachedValue() {
    if (this.error) {
      // invoke constructor of CachedValueResponse
      const resp = new CachedValueResponse(
        LookupStatus.Miss,
        OperationStatus.Error,
        this.error,
        null,
      );
      return resp;
    }
    const resp = new CachedValueResponse(
      ConvertCacheLookupStatus(this.checkResponse?.cachedValue?.lookupStatus),
      ConvertCacheOperationStatus(
        this.checkResponse?.cachedValue?.operationStatus,
      ),
      ConvertCacheError(this.checkResponse?.cachedValue?.error),
      this.checkResponse?.cachedValue?.value ?? null,
    );
    return resp;
  }

  /**
   * Gets the error associated with the flow.
   * @returns The error object.
   */
  Error() {
    return this.error;
  }

  /**
   * Gets the check response of the flow.
   * @returns The check response object.
   */
  CheckResponse() {
    return this.checkResponse;
  }

  /**
   * Gets the span associated with the flow.
   * @returns The span object.
   */
  Span() {
    return this.span;
  }

  /**
   * Ends the flow and performs necessary cleanup.
   */
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

      localCheckResponse.start = this.protoTimestampToJSON(
        this.checkResponse.start,
      );
      localCheckResponse.end = this.protoTimestampToJSON(
        this.checkResponse.end,
      );
      localCheckResponse.waitTime = this.protoDurationToJSON(
        this.checkResponse.waitTime,
      );

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

      this.span.setAttribute(
        CHECK_RESPONSE_LABEL,
        JSON.stringify(localCheckResponse),
      );
    }

    this.span.setAttribute(FLOW_STATUS_LABEL, this.status);
    this.span.setAttribute(FLOW_END_TIMESTAMP_LABEL, Date.now());
    this.span.end();
  }

  private protoTimestampToJSON(
    timestamp: _google_protobuf_Timestamp__Output | null,
  ) {
    if (timestamp) {
      return new Date(
        Number(timestamp.seconds) * 1000 + timestamp.nanos / 1000000,
      ).toISOString();
    }
    return timestamp;
  }

  private protoDurationToJSON(
    duration: _google_protobuf_Duration__Output | null,
  ) {
    if (duration) {
      return `${duration.seconds}.${duration.nanos}s`;
    }
    return duration;
  }
}
