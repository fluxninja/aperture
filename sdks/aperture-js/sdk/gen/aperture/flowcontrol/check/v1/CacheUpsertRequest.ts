// Original file: proto/flowcontrol/check/v1/check.proto

import type { CacheEntry as _aperture_flowcontrol_check_v1_CacheEntry, CacheEntry__Output as _aperture_flowcontrol_check_v1_CacheEntry__Output } from '../../../../aperture/flowcontrol/check/v1/CacheEntry';

export interface CacheUpsertRequest {
  'controlPoint'?: (string);
  'resultCacheEntry'?: (_aperture_flowcontrol_check_v1_CacheEntry | null);
  'globalCacheEntries'?: ({[key: string]: _aperture_flowcontrol_check_v1_CacheEntry});
}

export interface CacheUpsertRequest__Output {
  'controlPoint': (string);
  'resultCacheEntry': (_aperture_flowcontrol_check_v1_CacheEntry__Output | null);
  'globalCacheEntries': ({[key: string]: _aperture_flowcontrol_check_v1_CacheEntry__Output});
}
