// Original file: proto/flowcontrol/check/v1/check.proto

export const CacheLookupResult = {
  HIT: 0,
  MISS: 1,
} as const;

export type CacheLookupResult =
  | 'HIT'
  | 0
  | 'MISS'
  | 1

export type CacheLookupResult__Output = typeof CacheLookupResult[keyof typeof CacheLookupResult]
