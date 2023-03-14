from datetime import datetime, timezone
import enum
import time
from typing import Optional

from aperture_sdk._gen.aperture.flowcontrol.check.v1 import check_pb2
from aperture_sdk.const import (
    check_response_label,
    flow_end_timestamp_label,
    flow_status_label,
)
from google.protobuf import json_format
from opentelemetry import trace


class FlowStatus(enum.Enum):
    OK = enum.auto()
    Error = enum.auto()


class Flow:
    def __init__(
        self, span: trace.Span, check_response: Optional[check_pb2.CheckResponse]
    ):
        self._span = span
        self._check_response = check_response
        self._ended = False

    def accepted(self) -> bool:
        if not self._check_response:
            return True
        if (
            self._check_response.decision_type
            == check_pb2.CheckResponse.DECISION_TYPE_ACCEPTED
        ):
            return True
        return False

    def check_response(self) -> Optional[check_pb2.CheckResponse]:
        return self._check_response

    def end(self, status_code: FlowStatus) -> None:
        if self._ended:
            raise ValueError("flow already ended")
        self._ended = True

        check_response_json = (
            json_format.MessageToJson(self._check_response)
            if self._check_response
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
