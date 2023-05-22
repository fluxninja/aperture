import enum
import time
from contextlib import AbstractContextManager
from typing import Optional, TypeVar

from aperture_sdk._gen.aperture.flowcontrol.check.v1 import check_pb2
from aperture_sdk.const import (
    check_response_label,
    flow_end_timestamp_label,
    flow_status_label,
)
from google.protobuf import json_format
from opentelemetry import trace


class FlowResult(enum.Enum):
    Accepted = enum.auto()
    Rejected = enum.auto()
    Unreachable = enum.auto()


class FlowStatus(enum.Enum):
    OK = enum.auto()
    Error = enum.auto()


TFlow = TypeVar("TFlow", bound="Flow")


class Flow(AbstractContextManager):
    def __init__(
        self, span: trace.Span, check_response: Optional[check_pb2.CheckResponse]
    ):
        self._span = span
        self._check_response = check_response
        self._ended = False

    @property
    def result(self) -> FlowResult:
        if self.check_response is None:
            return FlowResult.Unreachable
        if (
            self.check_response.decision_type
            == check_pb2.CheckResponse.DECISION_TYPE_ACCEPTED
        ):
            return FlowResult.Accepted
        return FlowResult.Rejected

    @property
    def success(self) -> bool:
        return self.result != FlowResult.Unreachable

    @property
    def check_response(self) -> Optional[check_pb2.CheckResponse]:
        return self._check_response

    def end(self, status_code: FlowStatus) -> None:
        if self._ended:
            raise ValueError("flow already ended")
        self._ended = True

        check_response_json = (
            json_format.MessageToJson(self.check_response)
            if self.check_response
            else ""
        )
        self._span.set_attributes(
            {
                flow_status_label: status_code.name,
                check_response_label: check_response_json,
                flow_end_timestamp_label: time.monotonic_ns(),
            }
        )
        self._span.end()

    def __enter__(self: TFlow) -> TFlow:
        return self

    def __exit__(self, exc_type, _exc_value, _traceback) -> None:
        if self._ended:
            return
        if exc_type is not None:
            self.end(FlowStatus.Error)
        else:
            self.end(FlowStatus.OK)
