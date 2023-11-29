// Original file: proto/flowcontrol/check/v1/check.proto


export interface CacheDeleteRequest {
  'controlPoint'?: (string);
  'resultCacheKey'?: (string);
  'globalCacheKeys'?: (string)[];
}

export interface CacheDeleteRequest__Output {
  'controlPoint': (string);
  'resultCacheKey': (string);
  'globalCacheKeys': (string)[];
}
