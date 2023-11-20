// Original file: proto/flowcontrol/check/v1/check.proto

export const CacheResult = {
  Hit: 0,
  Miss: 1,
} as const;

export type CacheResult =
  | 'Hit'
  | 0
  | 'Miss'
  | 1

export type CacheResult__Output = typeof CacheResult[keyof typeof CacheResult]
