// Original file: proto/flowcontrol/check/v1/check.proto

export const CacheResponseCode = {
  SUCCESS: 0,
  ERROR: 1,
} as const;

export type CacheResponseCode =
  | 'SUCCESS'
  | 0
  | 'ERROR'
  | 1

export type CacheResponseCode__Output = typeof CacheResponseCode[keyof typeof CacheResponseCode]
