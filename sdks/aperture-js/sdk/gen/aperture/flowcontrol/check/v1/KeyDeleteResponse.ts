// Original file: proto/flowcontrol/check/v1/check.proto

import type { CacheOperationStatus as _aperture_flowcontrol_check_v1_CacheOperationStatus, CacheOperationStatus__Output as _aperture_flowcontrol_check_v1_CacheOperationStatus__Output } from '../../../../aperture/flowcontrol/check/v1/CacheOperationStatus';

export interface KeyDeleteResponse {
  'operationStatus'?: (_aperture_flowcontrol_check_v1_CacheOperationStatus);
  'error'?: (string);
}

export interface KeyDeleteResponse__Output {
  'operationStatus': (_aperture_flowcontrol_check_v1_CacheOperationStatus__Output);
  'error': (string);
}
