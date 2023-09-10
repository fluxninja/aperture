// Original file: proto/flowcontrol/check/v1/check.proto

import type { StatusCode as _aperture_flowcontrol_check_v1_StatusCode, StatusCode__Output as _aperture_flowcontrol_check_v1_StatusCode__Output } from '../../../../aperture/flowcontrol/check/v1/StatusCode';
import type { Long } from '@grpc/proto-loader';

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
  'schedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo | null);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo__Output {
  'label': (string);
  'schedulerInfo': (_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo__Output | null);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo {
  'remaining'?: (number | string);
  'current'?: (number | string);
  'label'?: (string);
  'tokensConsumed'?: (number | string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo__Output {
  'remaining': (number);
  'current': (number);
  'label': (string);
  'tokensConsumed': (number);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo {
  'label'?: (string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo__Output {
  'label': (string);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo {
  'workloadIndex'?: (string);
  'tokensConsumed'?: (number | string | Long);
}

export interface _aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo__Output {
  'workloadIndex': (string);
  'tokensConsumed': (string);
}

export interface LimiterDecision {
  'policyName'?: (string);
  'policyHash'?: (string);
  'componentId'?: (string);
  'dropped'?: (boolean);
  'reason'?: (_aperture_flowcontrol_check_v1_LimiterDecision_LimiterReason);
  'rateLimiterInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo | null);
  'loadSchedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo | null);
  'samplerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo | null);
  'quotaSchedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo | null);
  'deniedResponseStatusCode'?: (_aperture_flowcontrol_check_v1_StatusCode);
  'details'?: "rateLimiterInfo"|"loadSchedulerInfo"|"samplerInfo"|"quotaSchedulerInfo";
}

export interface LimiterDecision__Output {
  'policyName': (string);
  'policyHash': (string);
  'componentId': (string);
  'dropped': (boolean);
  'reason': (_aperture_flowcontrol_check_v1_LimiterDecision_LimiterReason__Output);
  'rateLimiterInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_RateLimiterInfo__Output | null);
  'loadSchedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_SchedulerInfo__Output | null);
  'samplerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_SamplerInfo__Output | null);
  'quotaSchedulerInfo'?: (_aperture_flowcontrol_check_v1_LimiterDecision_QuotaSchedulerInfo__Output | null);
  'deniedResponseStatusCode': (_aperture_flowcontrol_check_v1_StatusCode__Output);
}
