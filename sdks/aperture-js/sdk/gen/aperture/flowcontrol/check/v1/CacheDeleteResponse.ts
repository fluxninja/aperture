// Original file: proto/flowcontrol/check/v1/check.proto

import type { KeyDeleteResponse as _aperture_flowcontrol_check_v1_KeyDeleteResponse, KeyDeleteResponse__Output as _aperture_flowcontrol_check_v1_KeyDeleteResponse__Output } from '../../../../aperture/flowcontrol/check/v1/KeyDeleteResponse';

export interface CacheDeleteResponse {
  'resultCacheResponse'?: (_aperture_flowcontrol_check_v1_KeyDeleteResponse | null);
  'stateCacheResponses'?: ({[key: string]: _aperture_flowcontrol_check_v1_KeyDeleteResponse});
}

export interface CacheDeleteResponse__Output {
  'resultCacheResponse': (_aperture_flowcontrol_check_v1_KeyDeleteResponse__Output | null);
  'stateCacheResponses': ({[key: string]: _aperture_flowcontrol_check_v1_KeyDeleteResponse__Output});
}
