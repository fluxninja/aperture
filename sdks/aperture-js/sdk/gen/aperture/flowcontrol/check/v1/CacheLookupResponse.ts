// Original file: proto/flowcontrol/check/v1/check.proto

import type { KeyLookupResponse as _aperture_flowcontrol_check_v1_KeyLookupResponse, KeyLookupResponse__Output as _aperture_flowcontrol_check_v1_KeyLookupResponse__Output } from '../../../../aperture/flowcontrol/check/v1/KeyLookupResponse';

export interface CacheLookupResponse {
  'resultCacheResponse'?: (_aperture_flowcontrol_check_v1_KeyLookupResponse | null);
  'stateCacheResponses'?: ({[key: string]: _aperture_flowcontrol_check_v1_KeyLookupResponse});
}

export interface CacheLookupResponse__Output {
  'resultCacheResponse': (_aperture_flowcontrol_check_v1_KeyLookupResponse__Output | null);
  'stateCacheResponses': ({[key: string]: _aperture_flowcontrol_check_v1_KeyLookupResponse__Output});
}
