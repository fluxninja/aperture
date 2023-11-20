// Original file: proto/flowcontrol/check/v1/check.proto

export const CacheOperationStatus = {
  SUCCESS: 0,
  ERROR: 1,
} as const;

export type CacheOperationStatus =
  | 'SUCCESS'
  | 0
  | 'ERROR'
  | 1

export type CacheOperationStatus__Output = typeof CacheOperationStatus[keyof typeof CacheOperationStatus]
