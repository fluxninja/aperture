import { CacheLookupStatus } from "./gen/aperture/flowcontrol/check/v1/CacheLookupStatus.js";
import { CacheOperationStatus } from "./gen/aperture/flowcontrol/check/v1/CacheOperationStatus.js";

export enum LookupStatus {
  Hit = "HIT",
  Miss = "MISS",
}

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

export enum OperationStatus {
  Success = "SUCCESS",
  Error = "ERROR",
}

export function ConvertCacheOperationStatus(
  status: CacheOperationStatus | null | undefined,
): OperationStatus {
  switch (status) {
    case CacheOperationStatus.SUCCESS:
      return OperationStatus.Success;
    case CacheOperationStatus.ERROR:
      return OperationStatus.Error;
    default:
      return OperationStatus.Error;
  }
}

export function ConvertCacheError(
  error: string | null | undefined,
): Error | null {
  if (!error) {
    return null;
  }
  return new Error(error);
}

export class CachedValueResponse {
  lookupStatus: LookupStatus;
  operationStatus: OperationStatus;
  error: Error | null;
  value: Buffer | null;

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

  GetLookupStatus(): LookupStatus {
    return this.lookupStatus;
  }

  GetOperationStatus(): OperationStatus {
    return this.operationStatus;
  }

  GetError(): Error | null {
    return this.error;
  }

  GetValue(): Buffer | null {
    return this.value;
  }
}

export class SetCachedValueResponse {
  error: Error | null;
  operationStatus: OperationStatus;

  constructor(error: Error | null, operationStatus: OperationStatus) {
    this.error = error;
    this.operationStatus = operationStatus;
  }

  GetError(): Error | null {
    return this.error;
  }

  GetOperationStatus(): OperationStatus {
    return this.operationStatus;
  }
}

export class DeleteCachedValueResponse {
  error: Error | null;
  operationStatus: OperationStatus;

  constructor(error: Error | null, operationStatus: OperationStatus) {
    this.error = error;
    this.operationStatus = operationStatus;
  }

  GetError(): Error | null {
    return this.error;
  }

  GetOperationStatus(): OperationStatus {
    return this.operationStatus;
  }
}
