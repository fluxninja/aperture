// Original file: proto/flowcontrol/check/v1/check.proto


export interface CacheLookupRequest {
  'controlPoint'?: (string);
  'resultCacheKey'?: (string);
  'stateCacheKeys'?: (string)[];
}

export interface CacheLookupRequest__Output {
  'controlPoint': (string);
  'resultCacheKey': (string);
  'stateCacheKeys': (string)[];
}
