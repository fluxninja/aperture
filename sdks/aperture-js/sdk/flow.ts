import grpc from "@grpc/grpc-js";
import { Span } from "@opentelemetry/api";
import type {
  CacheEntry,
  KeyDeleteResponse,
  KeyLookupResponse,
  KeyUpsertResponse,
} from "./cache.js";
import { LookupStatus } from "./cache.js";
import {
  CHECK_RESPONSE_LABEL,
  FLOW_END_TIMESTAMP_LABEL,
  FLOW_START_TIMESTAMP_LABEL,
  FLOW_STATUS_LABEL,
  SOURCE_LABEL,
  WORKLOAD_START_TIMESTAMP_LABEL,
} from "./consts.js";
import type { CacheDeleteRequest } from "./gen/aperture/flowcontrol/check/v1/CacheDeleteRequest.js";
import { CacheLookupStatus } from "./gen/aperture/flowcontrol/check/v1/CacheLookupStatus.js";
import type { CacheUpsertRequest } from "./gen/aperture/flowcontrol/check/v1/CacheUpsertRequest.js";
import {
  CheckResponse__Output,
  _aperture_flowcontrol_check_v1_CheckResponse_DecisionType,
} from "./gen/aperture/flowcontrol/check/v1/CheckResponse.js";
import { FlowControlServiceClient } from "./gen/aperture/flowcontrol/check/v1/FlowControlService.js";
import { FlowEndRequest } from "./gen/aperture/flowcontrol/check/v1/FlowEndRequest.js";
import { FlowEndResponse as FlowEndResponseProto } from "./gen/aperture/flowcontrol/check/v1/FlowEndResponse.js";
import { InflightRequestRef } from "./gen/aperture/flowcontrol/check/v1/InflightRequestRef.js";
import type { Duration__Output as _google_protobuf_Duration__Output } from "./gen/google/protobuf/Duration";
import type { Timestamp__Output as _google_protobuf_Timestamp__Output } from "./gen/google/protobuf/Timestamp";

/**
 * Represents the status of a flow.
 */
export enum FlowStatus {
  OK = "OK",
  Error = "Error",
}

export interface Flow {
  checkResponse(): CheckResponse__Output | null;
  shouldRun(): boolean;
  setStatus(status: FlowStatus): void;
  // grpc options is optional argument
  setResultCache(
    cacheEntry: CacheEntry,
    grpcOptions?: grpc.CallOptions,
  ): Promise<KeyUpsertResponse>;
  setGlobalCache(
    key: string,
    cacheEntry: CacheEntry,
    grpcOptions?: grpc.CallOptions,
  ): Promise<KeyUpsertResponse>;
  deleteResultCache(
    grpcOptions?: grpc.CallOptions,
  ): Promise<KeyDeleteResponse | undefined>;
  deleteGlobalCache(
    key: string,
    grpcOptions?: grpc.CallOptions,
  ): Promise<KeyDeleteResponse>;
  resultCache(): KeyLookupResponse;
  globalCache(key: string): KeyLookupResponse;
  error(): Error | null;
  span(): Span;
  end(grpcOptions?: grpc.CallOptions): Promise<FlowEndResponse>;
  httpResponseCode(): Number | undefined;
  retryAfter(): { seconds: string | undefined; nanos: number | undefined };
}

/**
 * Represents a flow in the SDK.
 */
export class _Flow implements Flow {
  private ended: boolean = false;
  private status: FlowStatus = FlowStatus.OK;

  constructor(
    private fcsClient: FlowControlServiceClient,
    private controlPoint: string,
    private _span: Span,
    private startDate: number,
    private rampMode: boolean = false,
    private resultCacheKey: string | null = null,
    private _checkResponse: CheckResponse__Output | null = null,
    private _error: Error | null = null,
  ) {
    _span.setAttribute(SOURCE_LABEL, "sdk");
    _span.setAttribute(FLOW_START_TIMESTAMP_LABEL, this.startDate);
    _span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Date.now());
  }

  /**
   * Determines whether the flow should run based on the check response and ramp mode.
   * @returns A boolean value indicating whether the flow should run.
   */
  shouldRun() {
    if (
      (!this.rampMode && this._checkResponse === null) ||
      this._checkResponse?.decisionType ===
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
  setStatus(status: FlowStatus) {
    this.status = status;
  }

  /**
   * Sets the result cache entry for the flow.
   * @param cacheEntry The cache entry to set.
   * @returns A promise that resolves to the response of the key upsert operation.
   */
  async setResultCache(
    cacheEntry: CacheEntry,
    grpcCallOptions?: grpc.CallOptions,
  ) {
    return new Promise<KeyUpsertResponse>((resolve) => {
      if (!this.resultCacheKey) {
        const resp = new _KeyUpsertResponse(new Error("No cache key"));
        resolve(resp);
        return;
      }
      const cacheUpsertRequest: CacheUpsertRequest = {
        controlPoint: this.controlPoint,
        resultCacheEntry: {
          key: this.resultCacheKey,
          value: cacheEntry.value,
          ttl: cacheEntry.ttl,
        },
      };
      this.fcsClient.CacheUpsert(
        cacheUpsertRequest,
        grpcCallOptions ?? {},
        (err, res) => {
          if (err) {
            const resp = new _KeyUpsertResponse(err);
            resolve(resp);
            return;
          }
          if (!res) {
            const resp = new _KeyUpsertResponse(
              new Error("No cache upsert response"),
            );
            resolve(resp);
            return;
          }
          const resp = new _KeyUpsertResponse(
            convertError(res.resultCacheResponse?.error),
          );
          resolve(resp);
        },
      );
    });
  }

  /**
   * Sets a global cache entry for the flow.
   * @param key The key of the cache entry to set.
   * @param cacheEntry The cache entry to set.
   * @returns A promise that resolves to the response of the key upsert operation.
   */
  async setGlobalCache(
    key: string,
    cacheEntry: CacheEntry,
    grpcCallOptions?: grpc.CallOptions,
  ) {
    return new Promise<KeyUpsertResponse>((resolve) => {
      const cacheUpsertRequest: CacheUpsertRequest = {
        globalCacheEntries: {
          [key]: {
            value: cacheEntry.value,
            ttl: cacheEntry.ttl,
          },
        },
      };
      this.fcsClient.CacheUpsert(
        cacheUpsertRequest,
        grpcCallOptions ?? {},
        (err, res) => {
          if (err) {
            const resp = new _KeyUpsertResponse(err);
            resolve(resp);
            return;
          }
          if (!res) {
            const resp = new _KeyUpsertResponse(
              new Error("No cache upsert response"),
            );
            resolve(resp);
            return;
          }
          const resp = new _KeyUpsertResponse(
            convertError(res.globalCacheResponses[key]?.error),
          );
          resolve(resp);
        },
      );
    });
  }

  /**
   * Deletes the result cache for the flow.
   * @returns A promise that resolves to the response of the key delete operation.
   */
  async deleteResultCache(grpcCallOptions?: grpc.CallOptions) {
    if (!this.resultCacheKey) {
      return Promise.reject(new Error("No cache key"));
    }

    const key = this.resultCacheKey;
    return new Promise<KeyDeleteResponse | undefined>((resolve, _) => {
      const cacheDeleteRequest: CacheDeleteRequest = {
        controlPoint: this.controlPoint,
        resultCacheKey: key,
      };
      this.fcsClient.CacheDelete(
        cacheDeleteRequest,
        grpcCallOptions ?? {},
        (err, res) => {
          if (err) {
            const resp = new _KeyDeleteResponse(err);
            resolve(resp);
            return;
          }
          const resp = new _KeyDeleteResponse(
            convertError(res?.resultCacheResponse?.error),
          );
          resolve(resp);
        },
      );
    });
  }

  /**
   * Deletes a global cache entry for the flow.
   * @param key The key of the cache entry to delete.
   * @returns A promise that resolves to the response of the key delete operation.
   */
  async deleteGlobalCache(key: string, grpcCallOptions?: grpc.CallOptions) {
    return new Promise<KeyDeleteResponse>((resolve) => {
      const cacheDeleteRequest: CacheDeleteRequest = {
        globalCacheKeys: [key],
      };
      this.fcsClient.CacheDelete(
        cacheDeleteRequest,
        grpcCallOptions ?? {},
        (err, res) => {
          if (err) {
            const resp = new _KeyDeleteResponse(err);
            resolve(resp);
            return;
          }
          const resp = new _KeyDeleteResponse(
            convertError(res?.globalCacheResponses[key]?.error),
          );
          resolve(resp);
        },
      );
    });
  }

  /**
   * Returns result cache lookup response that was fetched at flow start.
   * @returns The result cache lookup response.
   */
  resultCache() {
    if (this._error) {
      // invoke constructor of CachedValueResponse
      const resp = new _KeyLookupResponse(LookupStatus.Miss, this._error, null);
      return resp;
    }

    if (!this._checkResponse) {
      const resp = new _KeyLookupResponse(
        LookupStatus.Miss,
        new Error("Check response was nil"),
        null,
      );
      return resp;
    }

    if (!this.shouldRun()) {
      const resp = new _KeyLookupResponse(
        LookupStatus.Miss,
        new Error("Flow was rejected"),
        null,
      );
      return resp;
    }

    const resultCacheResponse =
      this._checkResponse?.cacheLookupResponse?.resultCacheResponse;
    if (!resultCacheResponse) {
      // invoke constructor of CachedValueResponse
      const resp = new _KeyLookupResponse(
        LookupStatus.Miss,
        new Error("No result cache response found"),
        null,
      );
      return resp;
    }
    const resp = new _KeyLookupResponse(
      convertCacheLookupStatus(resultCacheResponse?.lookupStatus),
      convertError(resultCacheResponse?.error),
      resultCacheResponse?.value ?? null,
    );
    return resp;
  }

  /**
   * Returns global cache lookup response that was fetched at flow start.
   * @returns The global cache lookup response.
   */
  globalCache(key: string): KeyLookupResponse {
    if (this._error) {
      // invoke constructor of CachedValueResponse
      const resp = new _KeyLookupResponse(LookupStatus.Miss, this._error, null);
      return resp;
    }

    if (!this._checkResponse) {
      const resp = new _KeyLookupResponse(
        LookupStatus.Miss,
        new Error("Check response was nil"),
        null,
      );
      return resp;
    }

    if (!this.shouldRun()) {
      const resp = new _KeyLookupResponse(
        LookupStatus.Miss,
        new Error("Flow was rejected"),
        null,
      );
      return resp;
    }

    if (!this._checkResponse?.cacheLookupResponse?.globalCacheResponses) {
      // invoke constructor of CachedValueResponse
      const resp = new _KeyLookupResponse(
        LookupStatus.Miss,
        new Error("No global cache response found"),
        null,
      );
      return resp;
    }
    // if key is not found in global cache dict, return miss
    if (
      !this._checkResponse?.cacheLookupResponse?.globalCacheResponses?.hasOwnProperty(
        key,
      )
    ) {
      const resp = new _KeyLookupResponse(
        LookupStatus.Miss,
        new Error("Unknown global cache key"),
        null,
      );
      return resp;
    }

    const lookupResp =
      this._checkResponse?.cacheLookupResponse?.globalCacheResponses?.[key];
    const resp = new _KeyLookupResponse(
      convertCacheLookupStatus(lookupResp?.lookupStatus),
      convertError(lookupResp?.error),
      lookupResp?.value.byteLength ? lookupResp?.value : null,
    );

    return resp;
  }

  /**
   * Gets the error associated with the flow.
   * @returns The error object.
   */
  error() {
    return this._error;
  }

  /**
   * Gets the check response of the flow.
   * @returns The check response object.
   */
  checkResponse() {
    return this._checkResponse;
  }

  /**
   * Gets the retry-after duration of the flow.
   * @returns The retry-after duration.
   */
  retryAfter() {
    if (this._checkResponse) {
      return {
        seconds: this._checkResponse?.waitTime?.seconds.toString(),
        nanos: this._checkResponse?.waitTime?.nanos,
      };
    }
    return { seconds: undefined, nanos: undefined };
  }

  /**
   * Gets the status for the flow.
   * @returns The http response code.
   */
  httpResponseCode() {
    if (this._checkResponse) {
      let statusCode = this._checkResponse?.deniedResponseStatusCode;
      if (statusCode == 0) return 200;
      return Number(this._checkResponse?.deniedResponseStatusCode);
    }
    return 0;
  }

  /**
   * Gets the span associated with the flow.
   * @returns The span object.
   */
  span() {
    return this._span;
  }

  /**
   * Ends the flow and performs necessary cleanup.
   * @returns A promise that resolves to the response of the flow end operation. In the case of no ConcurrencyLimiter or ConcurrencyScheduler in limiter decisions, the promise resolves immediately.
   */
  async end(grpcCallOptions?: grpc.CallOptions) {
    return new Promise<FlowEndResponse>((resolve) => {
      if (this.ended) {
        return;
      }
      this.ended = true;

      var inflightRequestRefs = new Array<InflightRequestRef>();

      if (this._checkResponse) {
        // HACK: Change timestamps to ISO strings since the protobufjs library uses it in a different format
        // Issue: https://github.com/protobufjs/protobuf.js/issues/893
        // PR: https://github.com/protobufjs/protobuf.js/pull/1258
        // Current timestamp type: https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/timestamp.proto
        const localCheckResponse = this._checkResponse as any;

        if (
          this._checkResponse.cacheLookupResponse?.resultCacheResponse?.value
        ) {
          localCheckResponse.cacheLookupResponse.resultCacheResponse.value =
            bufferToByteArrayJson(
              this._checkResponse.cacheLookupResponse.resultCacheResponse.value,
            );
        }

        if (this._checkResponse.cacheLookupResponse?.globalCacheResponses) {
          Object.entries(
            this._checkResponse.cacheLookupResponse.globalCacheResponses,
          ).forEach(([key, value]) => {
            if (value.value) {
              localCheckResponse.cacheLookupResponse.globalCacheResponses[
                key
              ].value = bufferToByteArrayJson(value.value);
            }
          });
        }

        localCheckResponse.start = protoTimestampToJSON(
          this._checkResponse.start,
        );
        localCheckResponse.end = protoTimestampToJSON(this._checkResponse.end);
        localCheckResponse.waitTime = protoDurationToJSON(
          this._checkResponse.waitTime,
        );

        // Walk through individual decisions and convert waitTime fields,
        // then add to localCheckResponse, preserving immutability.
        if (this._checkResponse.limiterDecisions) {
          const decisions = this._checkResponse.limiterDecisions.map(
            (decision) => {
              return {
                ...decision,
                waitTime: protoDurationToJSON(decision.waitTime),
              };
            },
          );
          localCheckResponse.limiterDecisions = decisions;

          this._checkResponse.limiterDecisions.forEach((decision) => {
            if (decision.concurrencyLimiterInfo) {
              inflightRequestRefs.push({
                policyName: decision.policyName,
                policyHash: decision.policyHash,
                componentId: decision.componentId,
                label: decision.concurrencyLimiterInfo.label,
                requestId: decision.concurrencyLimiterInfo.requestId,
                tokens: decision.concurrencyLimiterInfo.tokensInfo?.consumed,
              });
            }
            if (decision.concurrencySchedulerInfo) {
              inflightRequestRefs.push({
                policyName: decision.policyName,
                policyHash: decision.policyHash,
                componentId: decision.componentId,
                label: decision.concurrencySchedulerInfo.label,
                requestId: decision.concurrencySchedulerInfo.requestId,
                tokens: decision.concurrencySchedulerInfo.tokensInfo?.consumed,
              });
            }
          });
        }

        this._span.setAttribute(
          CHECK_RESPONSE_LABEL,
          JSON.stringify(localCheckResponse),
        );
      }

      this._span.setAttribute(FLOW_STATUS_LABEL, this.status);
      this._span.setAttribute(FLOW_END_TIMESTAMP_LABEL, Date.now());
      this._span.end();

      const flowEndRequest: FlowEndRequest = {
        controlPoint: this.controlPoint,
        inflightRequests: inflightRequestRefs,
      };
      if (inflightRequestRefs.length > 0) {
        this.fcsClient.FlowEnd(
          flowEndRequest,
          grpcCallOptions ?? {},
          (err, res) => {
            if (err) {
              const resp = new _FlowEndResponse(err, {});
              resolve(resp);
              return;
            }
            if (!res) {
              const resp = new _FlowEndResponse(
                new Error("No flow end response"),
                {},
              );
              resolve(resp);
              return;
            }
            const resp = new _FlowEndResponse(null, res);
            resolve(resp);
            return;
          },
        );
      } else {
        if (!this._checkResponse) {
          const resp = new _FlowEndResponse(
            new Error("Check response was nil"),
            {},
          );
          resolve(resp);
          return;
        }
        const resp = new _FlowEndResponse(null, {});
        resolve(resp);
        return;
      }
    });
  }
}

/**
 * Represents a flow end response.
 */
interface FlowEndResponse {
  /**
   * Gets the response.
   * @returns FlowEndResponseProto.
   */
  getResponse(): FlowEndResponseProto;

  /**
   * Gets the error, if any.
   * @returns The error, or null if no error occurred.
   */
  getError(): Error | null;
}

function bufferToByteArrayJson(buffer: Buffer) {
  const byteArray = buffer.toJSON().data;
  var hash = "";
  for (var i = 0; i < byteArray.length; i++) {
    hash += String.fromCharCode(byteArray[i]);
  }
  return Buffer.from(hash).toString("base64");
}

function protoTimestampToJSON(
  timestamp: _google_protobuf_Timestamp__Output | null,
) {
  if (timestamp) {
    return new Date(
      Number(timestamp.seconds) * 1000 + timestamp.nanos / 1000000,
    ).toISOString();
  }
  return timestamp;
}

function protoDurationToJSON(
  duration: _google_protobuf_Duration__Output | null,
) {
  if (duration) {
    return `${duration.seconds}.${duration.nanos}s`;
  }
  return duration;
}

/**
 * Converts the cache lookup status to a lookup status.
 * @param status - The cache lookup status to convert.
 * @returns The converted lookup status.
 */
function convertCacheLookupStatus(
  status: CacheLookupStatus | null | undefined,
): LookupStatus {
  switch (status) {
    case CacheLookupStatus.HIT:
      return LookupStatus.Hit;
    case CacheLookupStatus.MISS:
      return LookupStatus.Miss;
    default:
      return LookupStatus.Miss;
  }
}

/**
 * Converts a cache error string into an Error object.
 *
 * @param error - The cache error string.
 * @returns The Error object representing the cache error, or null if the error string is empty.
 */
function convertError(error: string | undefined): Error | null {
  if (!error) {
    return null;
  }
  return new Error(error);
}

/**
 * Represents a cache value lookup.
 */
class _KeyLookupResponse {
  lookupStatus: LookupStatus;
  error: Error | null;
  value: Buffer | null;

  /**
   * Creates a new CachedValueResponse instance.
   * @param lookupStatus The lookup status.
   * @param error The error, if any.
   * @param value The cached value, if any.
   */
  constructor(
    lookupStatus: LookupStatus,
    error: Error | null,
    value: Buffer | null,
  ) {
    this.lookupStatus = lookupStatus;
    this.error = error;
    this.value = value;
  }

  /**
   * Gets the lookup status.
   * @returns The lookup status.
   */
  getLookupStatus(): LookupStatus {
    return this.lookupStatus;
  }

  /**
   * Gets the error, if any.
   * @returns The error, or null if no error occurred.
   */
  getError(): Error | null {
    return this.error;
  }

  /**
   * Gets the cached value, if any.
   * @returns The cached value, or null if no value is available.
   */
  getValue(): Buffer | null {
    return this.value;
  }
}

/**
 * Represents the response of updating or inserting a cache key.
 */
class _KeyUpsertResponse {
  private error: Error | null;

  /**
   * Creates a new instance of SetCachedValueResponse.
   * @param error The error that occurred during the operation, if any.
   */
  constructor(error: Error | null) {
    this.error = error;
  }

  /**
   * Gets the error that occurred during the operation.
   * @returns The error that occurred during the operation, or null if no error occurred.
   */
  getError(): Error | null {
    return this.error;
  }
}

/**
 * Represents the response of deleting a cache key.
 */
class _KeyDeleteResponse {
  private error: Error | null;

  /**
   * Creates a new instance of DeleteCachedValueResponse.
   * @param error The error that occurred during the delete operation, if any.
   */
  constructor(error: Error | null) {
    this.error = error;
  }

  /**
   * Gets the error that occurred during the delete operation, if any.
   * @returns The error object or null if no error occurred.
   */
  getError(): Error | null {
    return this.error;
  }
}

class _FlowEndResponse implements FlowEndResponse {
  private error: Error | null;
  private response: FlowEndResponseProto;

  /**
   * Creates a new instance of FlowEndResponse.
   * @param error The error that occurred during the operation, if any.
   */
  constructor(error: Error | null, response: FlowEndResponseProto) {
    this.error = error;
    this.response = response;
  }

  /**
   * Gets the response.
   * @returns FlowEndResponseProto.
   */
  getResponse(): FlowEndResponseProto {
    return this.response;
  }

  /**
   * Gets the error that occurred during the operation, if any.
   * @returns The error object or null if no error occurred.
   */
  getError(): Error | null {
    return this.error;
  }
}
