// Original file: proto/flowcontrol/check/v1/check.proto


export interface InflightRequestRef {
  'policyName'?: (string);
  'policyHash'?: (string);
  'componentId'?: (string);
  'label'?: (string);
  'requestId'?: (string);
  'tokens'?: (number | string);
}

export interface InflightRequestRef__Output {
  'policyName': (string);
  'policyHash': (string);
  'componentId': (string);
  'label': (string);
  'requestId': (string);
  'tokens': (number);
}
