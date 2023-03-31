from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class CheckRequest(_message.Message):
    __slots__ = ["control_point", "labels", "tokens"]
    class LabelsEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    CONTROL_POINT_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    TOKENS_FIELD_NUMBER: _ClassVar[int]
    control_point: str
    labels: _containers.ScalarMap[str, str]
    tokens: int
    def __init__(self, control_point: _Optional[str] = ..., labels: _Optional[_Mapping[str, str]] = ..., tokens: _Optional[int] = ...) -> None: ...

class CheckResponse(_message.Message):
    __slots__ = ["start", "end", "services", "control_point", "flow_label_keys", "telemetry_flow_labels", "decision_type", "reject_reason", "classifier_infos", "flux_meter_infos", "limiter_decisions"]
    class RejectReason(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = []
        REJECT_REASON_NONE: _ClassVar[CheckResponse.RejectReason]
        REJECT_REASON_RATE_LIMITED: _ClassVar[CheckResponse.RejectReason]
        REJECT_REASON_CONCURRENCY_LIMITED: _ClassVar[CheckResponse.RejectReason]
    REJECT_REASON_NONE: CheckResponse.RejectReason
    REJECT_REASON_RATE_LIMITED: CheckResponse.RejectReason
    REJECT_REASON_CONCURRENCY_LIMITED: CheckResponse.RejectReason
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
    def __init__(self, start: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., end: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., services: _Optional[_Iterable[str]] = ..., control_point: _Optional[str] = ..., flow_label_keys: _Optional[_Iterable[str]] = ..., telemetry_flow_labels: _Optional[_Mapping[str, str]] = ..., decision_type: _Optional[_Union[CheckResponse.DecisionType, str]] = ..., reject_reason: _Optional[_Union[CheckResponse.RejectReason, str]] = ..., classifier_infos: _Optional[_Iterable[_Union[ClassifierInfo, _Mapping]]] = ..., flux_meter_infos: _Optional[_Iterable[_Union[FluxMeterInfo, _Mapping]]] = ..., limiter_decisions: _Optional[_Iterable[_Union[LimiterDecision, _Mapping]]] = ...) -> None: ...

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
    __slots__ = ["policy_name", "policy_hash", "component_id", "dropped", "reason", "rate_limiter_info", "concurrency_limiter_info"]
    class LimiterReason(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = []
        LIMITER_REASON_UNSPECIFIED: _ClassVar[LimiterDecision.LimiterReason]
        LIMITER_REASON_KEY_NOT_FOUND: _ClassVar[LimiterDecision.LimiterReason]
    LIMITER_REASON_UNSPECIFIED: LimiterDecision.LimiterReason
    LIMITER_REASON_KEY_NOT_FOUND: LimiterDecision.LimiterReason
    class RateLimiterInfo(_message.Message):
        __slots__ = ["remaining", "current", "label"]
        REMAINING_FIELD_NUMBER: _ClassVar[int]
        CURRENT_FIELD_NUMBER: _ClassVar[int]
        LABEL_FIELD_NUMBER: _ClassVar[int]
        remaining: int
        current: int
        label: str
        def __init__(self, remaining: _Optional[int] = ..., current: _Optional[int] = ..., label: _Optional[str] = ...) -> None: ...
    class ConcurrencyLimiterInfo(_message.Message):
        __slots__ = ["workload_index"]
        WORKLOAD_INDEX_FIELD_NUMBER: _ClassVar[int]
        workload_index: str
        def __init__(self, workload_index: _Optional[str] = ...) -> None: ...
    POLICY_NAME_FIELD_NUMBER: _ClassVar[int]
    POLICY_HASH_FIELD_NUMBER: _ClassVar[int]
    COMPONENT_ID_FIELD_NUMBER: _ClassVar[int]
    DROPPED_FIELD_NUMBER: _ClassVar[int]
    REASON_FIELD_NUMBER: _ClassVar[int]
    RATE_LIMITER_INFO_FIELD_NUMBER: _ClassVar[int]
    CONCURRENCY_LIMITER_INFO_FIELD_NUMBER: _ClassVar[int]
    policy_name: str
    policy_hash: str
    component_id: str
    dropped: bool
    reason: LimiterDecision.LimiterReason
    rate_limiter_info: LimiterDecision.RateLimiterInfo
    concurrency_limiter_info: LimiterDecision.ConcurrencyLimiterInfo
    def __init__(self, policy_name: _Optional[str] = ..., policy_hash: _Optional[str] = ..., component_id: _Optional[str] = ..., dropped: bool = ..., reason: _Optional[_Union[LimiterDecision.LimiterReason, str]] = ..., rate_limiter_info: _Optional[_Union[LimiterDecision.RateLimiterInfo, _Mapping]] = ..., concurrency_limiter_info: _Optional[_Union[LimiterDecision.ConcurrencyLimiterInfo, _Mapping]] = ...) -> None: ...

class FluxMeterInfo(_message.Message):
    __slots__ = ["flux_meter_name"]
    FLUX_METER_NAME_FIELD_NUMBER: _ClassVar[int]
    flux_meter_name: str
    def __init__(self, flux_meter_name: _Optional[str] = ...) -> None: ...
