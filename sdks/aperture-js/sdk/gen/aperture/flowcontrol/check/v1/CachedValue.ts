// Original file: proto/flowcontrol/check/v1/check.proto

import type { CacheLookupResult as _aperture_flowcontrol_check_v1_CacheLookupResult, CacheLookupResult__Output as _aperture_flowcontrol_check_v1_CacheLookupResult__Output } from '../../../../aperture/flowcontrol/check/v1/CacheLookupResult';
import type { CacheResponseCode as _aperture_flowcontrol_check_v1_CacheResponseCode, CacheResponseCode__Output as _aperture_flowcontrol_check_v1_CacheResponseCode__Output } from '../../../../aperture/flowcontrol/check/v1/CacheResponseCode';

export interface CachedValue {
  'value'?: (Buffer | Uint8Array | string);
  'lookupResult'?: (_aperture_flowcontrol_check_v1_CacheLookupResult);
  'responseCode'?: (_aperture_flowcontrol_check_v1_CacheResponseCode);
  'message'?: (string);
}

export interface CachedValue__Output {
  'value': (Buffer);
  'lookupResult': (_aperture_flowcontrol_check_v1_CacheLookupResult__Output);
  'responseCode': (_aperture_flowcontrol_check_v1_CacheResponseCode__Output);
  'message': (string);
}
