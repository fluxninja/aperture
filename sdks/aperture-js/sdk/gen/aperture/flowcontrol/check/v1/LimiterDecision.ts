// Original file: proto/flowcontrol/check/v1/check.proto

import type { StatusCode as _aperture_flowcontrol_check_v1_StatusCode, StatusCode__Output as _aperture_flowcontrol_check_v1_StatusCode__Output } from '../../../../aperture/flowcontrol/check/v1/StatusCode';
import type { Duration as _google_protobuf_Duration, Duration__Output as _google_protobuf_Duration__Output } from '../../../../google/protobuf/Duration';

export interface _aperture_flowcontrol_check_v1_LimiterDecision_ConcurrencyLimiterInfo {
  'label'?: (string);
  'tokensInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo | null);
  'requestId'?: (string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_ConcurrencyLimiterInfo__Output {
  'label': (string);
  'tokensInfo': (_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo__Output | null);
  'requestId': (string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_ConcurrencySchedulerInfo {
  'label'?: (string);
  'workloadIndex'?: (string);
  'tokensInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo | null);
  'priority'?: (number | string);
  'requestId'?: (string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_ConcurrencySchedulerInfo__Output {
  'label': (string);
  'workloadIndex': (string);
  'tokensInfo': (_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo__Output | null);
  'priority': (number);
  'requestId': (string);
}

// Original file: proto/flowcontrol/check/v1/check.proto

export const _aperture_flowcontrol_check_v1_LimiterDecision_LimiterReason = {
  LIMITER_REASON_UNSPECIFIED: 0,
  LIMITER_REASON_KEY_NOT_FOUND: 1,
} as const;

export type _aperture_flowcontrol_check_v1_LimiterDecision_LimiterReason =
  | 'LIMITER_REASON_UNSPECIFIED'
  | 0
  | 'LIMITER_REASON_KEY_NOT_FOUND'
  | 1

export type _aperture_flowcontrol_check_v1_LimiterDecision_LimiterReason__Output = typeof _aperture_flowcontrol_check_v1_LimiterDecision_LimiterReason[keyof typeof _aperture_flowcontrol_check_v1_LimiterDecision_LimiterReason]

export interface _aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo {
  'label'?: (string);
  'workloadIndex'?: (string);
  'tokensInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo | null);
  'priority'?: (number | string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo__Output {
  'label': (string);
  'workloadIndex': (string);
  'tokensInfo': (_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo__Output | null);
  'priority': (number);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo {
  'label'?: (string);
  'tokensInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo | null);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo__Output {
  'label': (string);
  'tokensInfo': (_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo__Output | null);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo {
  'label'?: (string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo__Output {
  'label': (string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo {
  'workloadIndex'?: (string);
  'tokensInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo | null);
  'priority'?: (number | string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo__Output {
  'workloadIndex': (string);
  'tokensInfo': (_aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo__Output | null);
  'priority': (number);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo {
  'remaining'?: (number | string);
  'current'?: (number | string);
  'consumed'?: (number | string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_TokensInfo__Output {
  'remaining': (number);
  'current': (number);
  'consumed': (number);
}

export interface LimiterDecision {
  'policyName'?: (string);
  'policyHash'?: (string);
  'componentId'?: (string);
  'dropped'?: (boolean);
  'reason'?: (_aperture_flowcontrol_check_v1_LimiterDecision_LimiterReason);
  'deniedResponseStatusCode'?: (_aperture_flowcontrol_check_v1_StatusCode);
  'waitTime'?: (_google_protobuf_Duration | null);
  'rateLimiterInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo | null);
  'loadSchedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo | null);
  'samplerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo | null);
  'quotaSchedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo | null);
  'concurrencyLimiterInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_ConcurrencyLimiterInfo | null);
  'concurrencySchedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_ConcurrencySchedulerInfo | null);
  'details'?: "rateLimiterInfo"|"loadSchedulerInfo"|"samplerInfo"|"quotaSchedulerInfo"|"concurrencyLimiterInfo"|"concurrencySchedulerInfo";
}

export interface LimiterDecision__Output {
  'policyName': (string);
  'policyHash': (string);
  'componentId': (string);
  'dropped': (boolean);
  'reason': (_aperture_flowcontrol_check_v1_LimiterDecision_LimiterReason__Output);
  'deniedResponseStatusCode': (_aperture_flowcontrol_check_v1_StatusCode__Output);
  'waitTime': (_google_protobuf_Duration__Output | null);
  'rateLimiterInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo__Output | null);
  'loadSchedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo__Output | null);
  'samplerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo__Output | null);
  'quotaSchedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo__Output | null);
  'concurrencyLimiterInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_ConcurrencyLimiterInfo__Output | null);
  'concurrencySchedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_ConcurrencySchedulerInfo__Output | null);
}
