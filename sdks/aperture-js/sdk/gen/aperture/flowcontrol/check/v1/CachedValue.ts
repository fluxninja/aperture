// Original file: proto/flowcontrol/check/v1/check.proto

import type { LookupResult as _aperture_flowcontrol_check_v1_LookupResult, LookupResult__Output as _aperture_flowcontrol_check_v1_LookupResult__Output } from '../../../../aperture/flowcontrol/check/v1/LookupResult';

export interface CachedValue {
  'value'?: (Buffer | Uint8Array | string);
  'lookupResult'?: (_aperture_flowcontrol_check_v1_LookupResult);
}

export interface CachedValue__Output {
  'value': (Buffer);
  'lookupResult': (_aperture_flowcontrol_check_v1_LookupResult__Output);
}
