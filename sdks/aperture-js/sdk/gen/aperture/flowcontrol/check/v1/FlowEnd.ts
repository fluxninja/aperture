// Original file: proto/flowcontrol/check/v1/check.proto

import type { CacheItem as _aperture_flowcontrol_check_v1_CacheItem, CacheItem__Output as _aperture_flowcontrol_check_v1_CacheItem__Output } from '../../../../aperture/flowcontrol/check/v1/CacheItem';

export interface FlowEnd {
  'upserts'?: ({[key: string]: _aperture_flowcontrol_check_v1_CacheItem});
  'deletes'?: (string)[];
}

export interface FlowEnd__Output {
  'upserts': ({[key: string]: _aperture_flowcontrol_check_v1_CacheItem__Output});
  'deletes': (string)[];
}
