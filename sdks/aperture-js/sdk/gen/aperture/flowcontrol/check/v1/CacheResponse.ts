// Original file: proto/flowcontrol/check/v1/check.proto

import type { CacheResult as _aperture_flowcontrol_check_v1_CacheResult, CacheResult__Output as _aperture_flowcontrol_check_v1_CacheResult__Output } from '../../../../aperture/flowcontrol/check/v1/CacheResult';

export interface CacheResponse {
  'value'?: (Buffer | Uint8Array | string);
  'result'?: (_aperture_flowcontrol_check_v1_CacheResult);
}

export interface CacheResponse__Output {
  'value': (Buffer);
  'result': (_aperture_flowcontrol_check_v1_CacheResult__Output);
}
