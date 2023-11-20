// Original file: proto/flowcontrol/check/v1/check.proto

export const LookupResult = {
  Hit: 0,
  Miss: 1,
} as const;

export type LookupResult =
  | 'Hit'
  | 0
  | 'Miss'
  | 1

export type LookupResult__Output = typeof LookupResult[keyof typeof LookupResult]
