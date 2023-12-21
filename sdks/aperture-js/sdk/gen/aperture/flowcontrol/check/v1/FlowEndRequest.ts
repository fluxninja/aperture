// Original file: proto/flowcontrol/check/v1/check.proto

import type { InflightRequestRef as _aperture_flowcontrol_check_v1_InflightRequestRef, InflightRequestRef__Output as _aperture_flowcontrol_check_v1_InflightRequestRef__Output } from '../../../../aperture/flowcontrol/check/v1/InflightRequestRef';

export interface FlowEndRequest {
  'controlPoint'?: (string);
  'inflightRequests'?: (_aperture_flowcontrol_check_v1_InflightRequestRef)[];
}

export interface FlowEndRequest__Output {
  'controlPoint': (string);
  'inflightRequests': (_aperture_flowcontrol_check_v1_InflightRequestRef__Output)[];
}
