// Original file: proto/flowcontrol/check/v1/check.proto

import type { InflightRef as _aperture_flowcontrol_check_v1_InflightRef, InflightRef__Output as _aperture_flowcontrol_check_v1_InflightRef__Output } from '../../../../aperture/flowcontrol/check/v1/InflightRef';

export interface EndRequest {
  'inflightRequests'?: (_aperture_flowcontrol_check_v1_InflightRef)[];
}

export interface EndRequest__Output {
  'inflightRequests': (_aperture_flowcontrol_check_v1_InflightRef__Output)[];
}
