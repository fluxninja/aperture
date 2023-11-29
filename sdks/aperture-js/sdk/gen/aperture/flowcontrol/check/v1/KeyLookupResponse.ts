// Original file: proto/flowcontrol/check/v1/check.proto

import type { CacheLookupStatus as _aperture_flowcontrol_check_v1_CacheLookupStatus, CacheLookupStatus__Output as _aperture_flowcontrol_check_v1_CacheLookupStatus__Output } from '../../../../aperture/flowcontrol/check/v1/CacheLookupStatus';
import type { CacheOperationStatus as _aperture_flowcontrol_check_v1_CacheOperationStatus, CacheOperationStatus__Output as _aperture_flowcontrol_check_v1_CacheOperationStatus__Output } from '../../../../aperture/flowcontrol/check/v1/CacheOperationStatus';

export interface KeyLookupResponse {
  'value'?: (Buffer | Uint8Array | string);
  'lookupStatus'?: (_aperture_flowcontrol_check_v1_CacheLookupStatus);
  'operationStatus'?: (_aperture_flowcontrol_check_v1_CacheOperationStatus);
  'error'?: (string);
}

export interface KeyLookupResponse__Output {
  'value': (Buffer);
  'lookupStatus': (_aperture_flowcontrol_check_v1_CacheLookupStatus__Output);
  'operationStatus': (_aperture_flowcontrol_check_v1_CacheOperationStatus__Output);
  'error': (string);
}
