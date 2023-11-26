// Original file: proto/flowcontrol/check/v1/check.proto


export interface CheckRequest {
  'controlPoint'?: (string);
  'labels'?: ({[key: string]: string});
  'rampMode'?: (boolean);
  'resultCacheKey'?: (string);
  'stateCacheKeys'?: (string)[];
}

export interface CheckRequest__Output {
  'controlPoint': (string);
  'labels': ({[key: string]: string});
  'rampMode': (boolean);
  'resultCacheKey': (string);
  'stateCacheKeys': (string)[];
}
