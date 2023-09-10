// Original file: proto/flowcontrol/check/v1/check.proto

import type { Long } from '@grpc/proto-loader';

// Original file: proto/flowcontrol/check/v1/check.proto

export const _aperture_flowcontrol_check_v1_ClassifierInfo_Error = {
  ERROR_NONE: 0,
  ERROR_EVAL_FAILED: 1,
  ERROR_EMPTY_RESULTSET: 2,
  ERROR_AMBIGUOUS_RESULTSET: 3,
  ERROR_MULTI_EXPRESSION: 4,
  ERROR_EXPRESSION_NOT_MAP: 5,
} as const;

export type _aperture_flowcontrol_check_v1_ClassifierInfo_Error =
  | 'ERROR_NONE'
  | 0
  | 'ERROR_EVAL_FAILED'
  | 1
  | 'ERROR_EMPTY_RESULTSET'
  | 2
  | 'ERROR_AMBIGUOUS_RESULTSET'
  | 3
  | 'ERROR_MULTI_EXPRESSION'
  | 4
  | 'ERROR_EXPRESSION_NOT_MAP'
  | 5

export type _aperture_flowcontrol_check_v1_ClassifierInfo_Error__Output = typeof _aperture_flowcontrol_check_v1_ClassifierInfo_Error[keyof typeof _aperture_flowcontrol_check_v1_ClassifierInfo_Error]

export interface ClassifierInfo {
  'policyName'?: (string);
  'policyHash'?: (string);
  'classifierIndex'?: (number | string | Long);
  'error'?: (_aperture_flowcontrol_check_v1_ClassifierInfo_Error);
}

export interface ClassifierInfo__Output {
  'policyName': (string);
  'policyHash': (string);
  'classifierIndex': (Long);
  'error': (_aperture_flowcontrol_check_v1_ClassifierInfo_Error__Output);
}
