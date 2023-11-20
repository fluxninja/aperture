// Original file: proto/flowcontrol/check/v1/check.proto

export const CacheLookupStatus = {
  HIT: 0,
  MISS: 1,
} as const;

export type CacheLookupStatus =
  | 'HIT'
  | 0
  | 'MISS'
  | 1

export type CacheLookupStatus__Output = typeof CacheLookupStatus[keyof typeof CacheLookupStatus]
