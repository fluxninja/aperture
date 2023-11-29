// Original file: proto/flowcontrol/check/v1/check.proto

import type { Duration as _google_protobuf_Duration, Duration__Output as _google_protobuf_Duration__Output } from '../../../../google/protobuf/Duration';

export interface CacheEntry {
  'key'?: (string);
  'value'?: (Buffer | Uint8Array | string);
  'ttl'?: (_google_protobuf_Duration | null);
}

export interface CacheEntry__Output {
  'key': (string);
  'value': (Buffer);
  'ttl': (_google_protobuf_Duration__Output | null);
}
