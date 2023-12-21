// Original file: proto/flowcontrol/check/v1/check.proto


export interface InflightRef {
  'policyName'?: (string);
  'policyHash'?: (string);
  'componentId'?: (string);
  'label'?: (string);
  'requestId'?: (string);
  'tokens'?: (number | string);
  'ok'?: (boolean);
}

export interface InflightRef__Output {
  'policyName': (string);
  'policyHash': (string);
  'componentId': (string);
  'label': (string);
  'requestId': (string);
  'tokens': (number);
  'ok': (boolean);
}
