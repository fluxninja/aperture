import { CacheLookupStatus } from "./gen/aperture/flowcontrol/check/v1/CacheLookupStatus.js";
import { CacheOperationStatus } from "./gen/aperture/flowcontrol/check/v1/CacheOperationStatus.js";

/**
 * Represents the status of a lookup operation in the cache.
 */
export enum LookupStatus {
  Hit = "HIT",
  Miss = "MISS",
}

/**
 * Converts the cache lookup status to a lookup status.
 * @param status - The cache lookup status to convert.
 * @returns The converted lookup status.
 */
export function ConvertCacheLookupStatus(
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
 * Represents the status of an operation.
 */
export enum OperationStatus {
  Success = "SUCCESS",
  Error = "ERROR",
}

/**
 * Converts a cache operation status to an operation status.
 * @param status The cache operation status to convert.
 * @returns The converted operation status.
 */
export function ConvertCacheOperationStatus(status: CacheOperationStatus | undefined): OperationStatus {
  switch (status) {
    case CacheOperationStatus.SUCCESS:
      return OperationStatus.Success;
    case CacheOperationStatus.ERROR:
      return OperationStatus.Error;
    default:
      return OperationStatus.Error;
  }
}


/**
 * Converts a cache error string into an Error object.
 *
 * @param error - The cache error string.
 * @returns The Error object representing the cache error, or null if the error string is empty.
 */
export function ConvertCacheError(error: string | undefined): Error | null {
  if (!error) {
    return null;
  }
  return new Error(error);
}

/**
 * Represents a response from a cached value lookup.
 */
export class CachedValueResponse {
  lookupStatus: LookupStatus;
  operationStatus: OperationStatus;
  error: Error | null;
  value: Buffer | null;

  /**
   * Creates a new CachedValueResponse instance.
   * @param lookupStatus The lookup status.
   * @param operationStatus The operation status.
   * @param error The error, if any.
   * @param value The cached value, if any.
   */
  constructor(
    lookupStatus: LookupStatus,
    operationStatus: OperationStatus,
    error: Error | null,
    value: Buffer | null,
  ) {
    this.lookupStatus = lookupStatus;
    this.operationStatus = operationStatus;
    this.error = error;
    this.value = value;
  }

  /**
   * Gets the lookup status.
   * @returns The lookup status.
   */
  GetLookupStatus(): LookupStatus {
    return this.lookupStatus;
  }

  /**
   * Gets the operation status.
   * @returns The operation status.
   */
  GetOperationStatus(): OperationStatus {
    return this.operationStatus;
  }

  /**
   * Gets the error, if any.
   * @returns The error, or null if no error occurred.
   */
  GetError(): Error | null {
    return this.error;
  }

  /**
   * Gets the cached value, if any.
   * @returns The cached value, or null if no value is available.
   */
  GetValue(): Buffer | null {
    return this.value;
  }
}

/**
 * Represents the response of setting a cached value.
 */
export class SetCachedValueResponse {
  error: Error | null;
  operationStatus: OperationStatus;

  /**
   * Creates a new instance of SetCachedValueResponse.
   * @param error The error that occurred during the operation, if any.
   * @param operationStatus The status of the operation.
   */
  constructor(error: Error | null, operationStatus: OperationStatus) {
    this.error = error;
    this.operationStatus = operationStatus;
  }

  /**
   * Gets the error that occurred during the operation.
   * @returns The error that occurred during the operation, or null if no error occurred.
   */
  GetError(): Error | null {
    return this.error;
  }

  /**
   * Gets the status of the operation.
   * @returns The status of the operation.
   */
  GetOperationStatus(): OperationStatus {
    return this.operationStatus;
  }
}

/**
 * Represents the response of deleting a cached value.
 */
export class DeleteCachedValueResponse {
  error: Error | null;
  operationStatus: OperationStatus;

  /**
   * Creates a new instance of DeleteCachedValueResponse.
   * @param error The error that occurred during the delete operation, if any.
   * @param operationStatus The status of the delete operation.
   */
  constructor(error: Error | null, operationStatus: OperationStatus) {
    this.error = error;
    this.operationStatus = operationStatus;
  }

  /**
   * Gets the error that occurred during the delete operation, if any.
   * @returns The error object or null if no error occurred.
   */
  GetError(): Error | null {
    return this.error;
  }

  /**
   * Gets the status of the delete operation.
   * @returns The operation status.
   */
  GetOperationStatus(): OperationStatus {
    return this.operationStatus;
  }
}
