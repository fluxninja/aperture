from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class StatusCode(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    Empty: _ClassVar[StatusCode]
    Continue: _ClassVar[StatusCode]
    OK: _ClassVar[StatusCode]
    Created: _ClassVar[StatusCode]
    Accepted: _ClassVar[StatusCode]
    NonAuthoritativeInformation: _ClassVar[StatusCode]
    NoContent: _ClassVar[StatusCode]
    ResetContent: _ClassVar[StatusCode]
    PartialContent: _ClassVar[StatusCode]
    MultiStatus: _ClassVar[StatusCode]
    AlreadyReported: _ClassVar[StatusCode]
    IMUsed: _ClassVar[StatusCode]
    MultipleChoices: _ClassVar[StatusCode]
    MovedPermanently: _ClassVar[StatusCode]
    Found: _ClassVar[StatusCode]
    SeeOther: _ClassVar[StatusCode]
    NotModified: _ClassVar[StatusCode]
    UseProxy: _ClassVar[StatusCode]
    TemporaryRedirect: _ClassVar[StatusCode]
    PermanentRedirect: _ClassVar[StatusCode]
    BadRequest: _ClassVar[StatusCode]
    Unauthorized: _ClassVar[StatusCode]
    PaymentRequired: _ClassVar[StatusCode]
    Forbidden: _ClassVar[StatusCode]
    NotFound: _ClassVar[StatusCode]
    MethodNotAllowed: _ClassVar[StatusCode]
    NotAcceptable: _ClassVar[StatusCode]
    ProxyAuthenticationRequired: _ClassVar[StatusCode]
    RequestTimeout: _ClassVar[StatusCode]
    Conflict: _ClassVar[StatusCode]
    Gone: _ClassVar[StatusCode]
    LengthRequired: _ClassVar[StatusCode]
    PreconditionFailed: _ClassVar[StatusCode]
    PayloadTooLarge: _ClassVar[StatusCode]
    URITooLong: _ClassVar[StatusCode]
    UnsupportedMediaType: _ClassVar[StatusCode]
    RangeNotSatisfiable: _ClassVar[StatusCode]
    ExpectationFailed: _ClassVar[StatusCode]
    MisdirectedRequest: _ClassVar[StatusCode]
    UnprocessableEntity: _ClassVar[StatusCode]
    Locked: _ClassVar[StatusCode]
    FailedDependency: _ClassVar[StatusCode]
    UpgradeRequired: _ClassVar[StatusCode]
    PreconditionRequired: _ClassVar[StatusCode]
    TooManyRequests: _ClassVar[StatusCode]
    RequestHeaderFieldsTooLarge: _ClassVar[StatusCode]
    InternalServerError: _ClassVar[StatusCode]
    NotImplemented: _ClassVar[StatusCode]
    BadGateway: _ClassVar[StatusCode]
    ServiceUnavailable: _ClassVar[StatusCode]
    GatewayTimeout: _ClassVar[StatusCode]
    HTTPVersionNotSupported: _ClassVar[StatusCode]
    VariantAlsoNegotiates: _ClassVar[StatusCode]
    InsufficientStorage: _ClassVar[StatusCode]
    LoopDetected: _ClassVar[StatusCode]
    NotExtended: _ClassVar[StatusCode]
    NetworkAuthenticationRequired: _ClassVar[StatusCode]
Empty: StatusCode
Continue: StatusCode
OK: StatusCode
Created: StatusCode
Accepted: StatusCode
NonAuthoritativeInformation: StatusCode
NoContent: StatusCode
ResetContent: StatusCode
PartialContent: StatusCode
MultiStatus: StatusCode
AlreadyReported: StatusCode
IMUsed: StatusCode
MultipleChoices: StatusCode
MovedPermanently: StatusCode
Found: StatusCode
SeeOther: StatusCode
NotModified: StatusCode
UseProxy: StatusCode
TemporaryRedirect: StatusCode
PermanentRedirect: StatusCode
BadRequest: StatusCode
Unauthorized: StatusCode
PaymentRequired: StatusCode
Forbidden: StatusCode
NotFound: StatusCode
MethodNotAllowed: StatusCode
NotAcceptable: StatusCode
ProxyAuthenticationRequired: StatusCode
RequestTimeout: StatusCode
Conflict: StatusCode
Gone: StatusCode
LengthRequired: StatusCode
PreconditionFailed: StatusCode
PayloadTooLarge: StatusCode
URITooLong: StatusCode
UnsupportedMediaType: StatusCode
RangeNotSatisfiable: StatusCode
ExpectationFailed: StatusCode
MisdirectedRequest: StatusCode
UnprocessableEntity: StatusCode
Locked: StatusCode
FailedDependency: StatusCode
UpgradeRequired: StatusCode
PreconditionRequired: StatusCode
TooManyRequests: StatusCode
RequestHeaderFieldsTooLarge: StatusCode
InternalServerError: StatusCode
NotImplemented: StatusCode
BadGateway: StatusCode
ServiceUnavailable: StatusCode
GatewayTimeout: StatusCode
HTTPVersionNotSupported: StatusCode
VariantAlsoNegotiates: StatusCode
InsufficientStorage: StatusCode
LoopDetected: StatusCode
NotExtended: StatusCode
NetworkAuthenticationRequired: StatusCode

class CheckRequest(_message.Message):
    __slots__ = ["control_point", "labels", "ramp_mode"]
    class LabelsEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    CONTROL_POINT_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    RAMP_MODE_FIELD_NUMBER: _ClassVar[int]
    control_point: str
    labels: _containers.ScalarMap[str, str]
    ramp_mode: bool
    def __init__(self, control_point: _Optional[str] = ..., labels: _Optional[_Mapping[str, str]] = ..., ramp_mode: bool = ...) -> None: ...

class CheckResponse(_message.Message):
    __slots__ = ["start", "end", "services", "control_point", "flow_label_keys", "telemetry_flow_labels", "decision_type", "reject_reason", "classifier_infos", "flux_meter_infos", "limiter_decisions", "wait_time", "denied_response_status_code"]
    class RejectReason(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = []
        REJECT_REASON_NONE: _ClassVar[CheckResponse.RejectReason]
        REJECT_REASON_RATE_LIMITED: _ClassVar[CheckResponse.RejectReason]
        REJECT_REASON_NO_TOKENS: _ClassVar[CheckResponse.RejectReason]
        REJECT_REASON_NOT_SAMPLED: _ClassVar[CheckResponse.RejectReason]
        REJECT_REASON_NO_MATCHING_RAMP: _ClassVar[CheckResponse.RejectReason]
    REJECT_REASON_NONE: CheckResponse.RejectReason
    REJECT_REASON_RATE_LIMITED: CheckResponse.RejectReason
    REJECT_REASON_NO_TOKENS: CheckResponse.RejectReason
    REJECT_REASON_NOT_SAMPLED: CheckResponse.RejectReason
    REJECT_REASON_NO_MATCHING_RAMP: CheckResponse.RejectReason
    class DecisionType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = []
        DECISION_TYPE_ACCEPTED: _ClassVar[CheckResponse.DecisionType]
        DECISION_TYPE_REJECTED: _ClassVar[CheckResponse.DecisionType]
    DECISION_TYPE_ACCEPTED: CheckResponse.DecisionType
    DECISION_TYPE_REJECTED: CheckResponse.DecisionType
    class TelemetryFlowLabelsEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    SERVICES_FIELD_NUMBER: _ClassVar[int]
    CONTROL_POINT_FIELD_NUMBER: _ClassVar[int]
    FLOW_LABEL_KEYS_FIELD_NUMBER: _ClassVar[int]
    TELEMETRY_FLOW_LABELS_FIELD_NUMBER: _ClassVar[int]
    DECISION_TYPE_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    CLASSIFIER_INFOS_FIELD_NUMBER: _ClassVar[int]
    FLUX_METER_INFOS_FIELD_NUMBER: _ClassVar[int]
    LIMITER_DECISIONS_FIELD_NUMBER: _ClassVar[int]
    WAIT_TIME_FIELD_NUMBER: _ClassVar[int]
    DENIED_RESPONSE_STATUS_CODE_FIELD_NUMBER: _ClassVar[int]
    start: _timestamp_pb2.Timestamp
    end: _timestamp_pb2.Timestamp
    services: _containers.RepeatedScalarFieldContainer[str]
    control_point: str
    flow_label_keys: _containers.RepeatedScalarFieldContainer[str]
    telemetry_flow_labels: _containers.ScalarMap[str, str]
    decision_type: CheckResponse.DecisionType
    reject_reason: CheckResponse.RejectReason
    classifier_infos: _containers.RepeatedCompositeFieldContainer[ClassifierInfo]
    flux_meter_infos: _containers.RepeatedCompositeFieldContainer[FluxMeterInfo]
    limiter_decisions: _containers.RepeatedCompositeFieldContainer[LimiterDecision]
    wait_time: _duration_pb2.Duration
    denied_response_status_code: StatusCode
    def __init__(self, start: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., end: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., services: _Optional[_Iterable[str]] = ..., control_point: _Optional[str] = ..., flow_label_keys: _Optional[_Iterable[str]] = ..., telemetry_flow_labels: _Optional[_Mapping[str, str]] = ..., decision_type: _Optional[_Union[CheckResponse.DecisionType, str]] = ..., reject_reason: _Optional[_Union[CheckResponse.RejectReason, str]] = ..., classifier_infos: _Optional[_Iterable[_Union[ClassifierInfo, _Mapping]]] = ..., flux_meter_infos: _Optional[_Iterable[_Union[FluxMeterInfo, _Mapping]]] = ..., limiter_decisions: _Optional[_Iterable[_Union[LimiterDecision, _Mapping]]] = ..., wait_time: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., denied_response_status_code: _Optional[_Union[StatusCode, str]] = ...) -> None: ...

class ClassifierInfo(_message.Message):
    __slots__ = ["policy_name", "policy_hash", "classifier_index", "error"]
    class Error(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = []
        ERROR_NONE: _ClassVar[ClassifierInfo.Error]
        ERROR_EVAL_FAILED: _ClassVar[ClassifierInfo.Error]
        ERROR_EMPTY_RESULTSET: _ClassVar[ClassifierInfo.Error]
        ERROR_AMBIGUOUS_RESULTSET: _ClassVar[ClassifierInfo.Error]
        ERROR_MULTI_EXPRESSION: _ClassVar[ClassifierInfo.Error]
        ERROR_EXPRESSION_NOT_MAP: _ClassVar[ClassifierInfo.Error]
    ERROR_NONE: ClassifierInfo.Error
    ERROR_EVAL_FAILED: ClassifierInfo.Error
    ERROR_EMPTY_RESULTSET: ClassifierInfo.Error
    ERROR_AMBIGUOUS_RESULTSET: ClassifierInfo.Error
    ERROR_MULTI_EXPRESSION: ClassifierInfo.Error
    ERROR_EXPRESSION_NOT_MAP: ClassifierInfo.Error
    POLICY_NAME_FIELD_NUMBER: _ClassVar[int]
    POLICY_HASH_FIELD_NUMBER: _ClassVar[int]
    CLASSIFIER_INDEX_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    policy_name: str
    policy_hash: str
    classifier_index: int
    error: ClassifierInfo.Error
    def __init__(self, policy_name: _Optional[str] = ..., policy_hash: _Optional[str] = ..., classifier_index: _Optional[int] = ..., error: _Optional[_Union[ClassifierInfo.Error, str]] = ...) -> None: ...

class LimiterDecision(_message.Message):
    __slots__ = ["policy_name", "policy_hash", "component_id", "dropped", "reason", "denied_response_status_code", "wait_time", "rate_limiter_info", "load_scheduler_info", "sampler_info", "quota_scheduler_info"]
    class LimiterReason(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = []
        LIMITER_REASON_UNSPECIFIED: _ClassVar[LimiterDecision.LimiterReason]
        LIMITER_REASON_KEY_NOT_FOUND: _ClassVar[LimiterDecision.LimiterReason]
    LIMITER_REASON_UNSPECIFIED: LimiterDecision.LimiterReason
    LIMITER_REASON_KEY_NOT_FOUND: LimiterDecision.LimiterReason
    class TokensInfo(_message.Message):
        __slots__ = ["remaining", "current", "consumed"]
        REMAINING_FIELD_NUMBER: _ClassVar[int]
        CURRENT_FIELD_NUMBER: _ClassVar[int]
        CONSUMED_FIELD_NUMBER: _ClassVar[int]
        remaining: float
        current: float
        consumed: float
        def __init__(self, remaining: _Optional[float] = ..., current: _Optional[float] = ..., consumed: _Optional[float] = ...) -> None: ...
    class RateLimiterInfo(_message.Message):
        __slots__ = ["label", "tokens_info"]
        LABEL_FIELD_NUMBER: _ClassVar[int]
        TOKENS_INFO_FIELD_NUMBER: _ClassVar[int]
        label: str
        tokens_info: LimiterDecision.TokensInfo
        def __init__(self, label: _Optional[str] = ..., tokens_info: _Optional[_Union[LimiterDecision.TokensInfo, _Mapping]] = ...) -> None: ...
    class SchedulerInfo(_message.Message):
        __slots__ = ["workload_index", "tokens_info", "priority"]
        WORKLOAD_INDEX_FIELD_NUMBER: _ClassVar[int]
        TOKENS_INFO_FIELD_NUMBER: _ClassVar[int]
        PRIORITY_FIELD_NUMBER: _ClassVar[int]
        workload_index: str
        tokens_info: LimiterDecision.TokensInfo
        priority: float
        def __init__(self, workload_index: _Optional[str] = ..., tokens_info: _Optional[_Union[LimiterDecision.TokensInfo, _Mapping]] = ..., priority: _Optional[float] = ...) -> None: ...
    class SamplerInfo(_message.Message):
        __slots__ = ["label"]
        LABEL_FIELD_NUMBER: _ClassVar[int]
        label: str
        def __init__(self, label: _Optional[str] = ...) -> None: ...
    class QuotaSchedulerInfo(_message.Message):
        __slots__ = ["label", "workload_index", "tokens_info", "priority"]
        LABEL_FIELD_NUMBER: _ClassVar[int]
        WORKLOAD_INDEX_FIELD_NUMBER: _ClassVar[int]
        TOKENS_INFO_FIELD_NUMBER: _ClassVar[int]
        PRIORITY_FIELD_NUMBER: _ClassVar[int]
        label: str
        workload_index: str
        tokens_info: LimiterDecision.TokensInfo
        priority: float
        def __init__(self, label: _Optional[str] = ..., workload_index: _Optional[str] = ..., tokens_info: _Optional[_Union[LimiterDecision.TokensInfo, _Mapping]] = ..., priority: _Optional[float] = ...) -> None: ...
    POLICY_NAME_FIELD_NUMBER: _ClassVar[int]
    POLICY_HASH_FIELD_NUMBER: _ClassVar[int]
    COMPONENT_ID_FIELD_NUMBER: _ClassVar[int]
    DROPPED_FIELD_NUMBER: _ClassVar[int]
    REASON_FIELD_NUMBER: _ClassVar[int]
    DENIED_RESPONSE_STATUS_CODE_FIELD_NUMBER: _ClassVar[int]
    WAIT_TIME_FIELD_NUMBER: _ClassVar[int]
    RATE_LIMITER_INFO_FIELD_NUMBER: _ClassVar[int]
    LOAD_SCHEDULER_INFO_FIELD_NUMBER: _ClassVar[int]
    SAMPLER_INFO_FIELD_NUMBER: _ClassVar[int]
    QUOTA_SCHEDULER_INFO_FIELD_NUMBER: _ClassVar[int]
    policy_name: str
    policy_hash: str
    component_id: str
    dropped: bool
    reason: LimiterDecision.LimiterReason
    denied_response_status_code: StatusCode
    wait_time: _duration_pb2.Duration
    rate_limiter_info: LimiterDecision.RateLimiterInfo
    load_scheduler_info: LimiterDecision.SchedulerInfo
    sampler_info: LimiterDecision.SamplerInfo
    quota_scheduler_info: LimiterDecision.QuotaSchedulerInfo
    def __init__(self, policy_name: _Optional[str] = ..., policy_hash: _Optional[str] = ..., component_id: _Optional[str] = ..., dropped: bool = ..., reason: _Optional[_Union[LimiterDecision.LimiterReason, str]] = ..., denied_response_status_code: _Optional[_Union[StatusCode, str]] = ..., wait_time: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., rate_limiter_info: _Optional[_Union[LimiterDecision.RateLimiterInfo, _Mapping]] = ..., load_scheduler_info: _Optional[_Union[LimiterDecision.SchedulerInfo, _Mapping]] = ..., sampler_info: _Optional[_Union[LimiterDecision.SamplerInfo, _Mapping]] = ..., quota_scheduler_info: _Optional[_Union[LimiterDecision.QuotaSchedulerInfo, _Mapping]] = ...) -> None: ...

class FluxMeterInfo(_message.Message):
    __slots__ = ["flux_meter_name"]
    FLUX_METER_NAME_FIELD_NUMBER: _ClassVar[int]
    flux_meter_name: str
    def __init__(self, flux_meter_name: _Optional[str] = ...) -> None: ...
