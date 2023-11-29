// Original file: proto/flowcontrol/check/v1/check.proto

import type { KeyUpsertResponse as _aperture_flowcontrol_check_v1_KeyUpsertResponse, KeyUpsertResponse__Output as _aperture_flowcontrol_check_v1_KeyUpsertResponse__Output } from '../../../../aperture/flowcontrol/check/v1/KeyUpsertResponse';

export interface CacheUpsertResponse {
  'resultCacheResponse'?: (_aperture_flowcontrol_check_v1_KeyUpsertResponse | null);
  'globalCacheResponses'?: ({[key: string]: _aperture_flowcontrol_check_v1_KeyUpsertResponse});
}

export interface CacheUpsertResponse__Output {
  'resultCacheResponse': (_aperture_flowcontrol_check_v1_KeyUpsertResponse__Output | null);
  'globalCacheResponses': ({[key: string]: _aperture_flowcontrol_check_v1_KeyUpsertResponse__Output});
}
