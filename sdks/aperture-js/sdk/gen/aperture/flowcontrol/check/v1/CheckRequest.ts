// Original file: proto/flowcontrol/check/v1/check.proto

import type { CacheLookupRequest as _aperture_flowcontrol_check_v1_CacheLookupRequest, CacheLookupRequest__Output as _aperture_flowcontrol_check_v1_CacheLookupRequest__Output } from '../../../../aperture/flowcontrol/check/v1/CacheLookupRequest';

export interface CheckRequest {
  'controlPoint'?: (string);
  'labels'?: ({[key: string]: string});
  'rampMode'?: (boolean);
  'cacheLookupRequest'?: (_aperture_flowcontrol_check_v1_CacheLookupRequest | null);
  'expectEnd'?: (boolean);
}

export interface CheckRequest__Output {
  'controlPoint': (string);
  'labels': ({[key: string]: string});
  'rampMode': (boolean);
  'cacheLookupRequest': (_aperture_flowcontrol_check_v1_CacheLookupRequest__Output | null);
  'expectEnd': (boolean);
}
