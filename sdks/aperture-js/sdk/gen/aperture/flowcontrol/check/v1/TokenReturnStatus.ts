// Original file: proto/flowcontrol/check/v1/check.proto

import type { InflightRequestRef as _aperture_flowcontrol_check_v1_InflightRequestRef, InflightRequestRef__Output as _aperture_flowcontrol_check_v1_InflightRequestRef__Output } from '../../../../aperture/flowcontrol/check/v1/InflightRequestRef';

export interface TokenReturnStatus {
  'inflightRequestRef'?: (_aperture_flowcontrol_check_v1_InflightRequestRef | null);
  'returned'?: (boolean);
  'error'?: (string);
}

export interface TokenReturnStatus__Output {
  'inflightRequestRef': (_aperture_flowcontrol_check_v1_InflightRequestRef__Output | null);
  'returned': (boolean);
  'error': (string);
}
