// Original file: null

import type { Long } from '@grpc/proto-loader';

export interface Duration {
  'seconds'?: (number | string | Long);
  'nanos'?: (number);
}

export interface Duration__Output {
  'seconds': (Long);
  'nanos': (number);
}
