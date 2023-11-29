// Original file: proto/flowcontrol/check/v1/check.proto

import type { CacheOperationStatus as _aperture_flowcontrol_check_v1_CacheOperationStatus, CacheOperationStatus__Output as _aperture_flowcontrol_check_v1_CacheOperationStatus__Output } from '../../../../aperture/flowcontrol/check/v1/CacheOperationStatus';

export interface CacheEntryStatus {
  'key'?: (string);
  'operationStatus'?: (_aperture_flowcontrol_check_v1_CacheOperationStatus);
  'error'?: (string);
}

export interface CacheEntryStatus__Output {
  'key': (string);
  'operationStatus': (_aperture_flowcontrol_check_v1_CacheOperationStatus__Output);
  'error': (string);
}
