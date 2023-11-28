import { Duration } from "@grpc/grpc-js/build/src/duration.js";

/**
 * Represents the status of a lookup operation in the cache.
 */
export enum LookupStatus {
  Hit = "HIT",
  Miss = "MISS",
}

/**
 * Represents a cache value lookup.
 */
export interface KeyLookupResponse {
  /**
   * Gets the lookup status.
   * @returns The lookup status.
   */
  getLookupStatus(): LookupStatus;

  /**
   * Gets the error, if any.
   * @returns The error, or null if no error occurred.
   */
  getError(): Error | null;

  /**
   * Gets the cached value, if any.
   * @returns The cached value, or null if no value is available.
   */
  getValue(): Buffer | null;
}

/**
 * Represents a cache entry.
 */
export interface CacheEntry {
  value: Buffer;
  ttl: Duration;
}

/**
 * Represents the response of a cache key update operation.
 */
export interface KeyUpsertResponse {
  /**
   * Gets the error that occurred during the operation.
   * @returns The error that occurred during the operation, or null if no error occurred.
   */
  getError(): Error | null;
}

/**
 * Represents the response of a cache key delete operation.
 */
export interface KeyDeleteResponse {
  /**
   * Gets the error that occurred during the delete operation, if any.
   * @returns The error object or null if no error occurred.
   */
  getError(): Error | null;
}
