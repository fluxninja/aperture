import datetime
import enum
import logging
import time
from contextlib import AbstractContextManager
from typing import Optional, TypeVar

from aperture_sdk._gen.aperture.flowcontrol.check.v1 import check_pb2
from aperture_sdk._gen.aperture.flowcontrol.check.v1.check_pb2 import (
    MISS,
    CacheDeleteResponse,
    CacheUpsertResponse,
)
from aperture_sdk._gen.aperture.flowcontrol.check.v1.check_pb2_grpc import (
    FlowControlServiceStub,
)
from aperture_sdk.cache import *
from aperture_sdk.const import (
    check_response_label,
    flow_end_timestamp_label,
    flow_status_label,
)
from google.protobuf import json_format
from opentelemetry import trace


class FlowDecision(enum.Enum):
    Accepted = enum.auto()
    Rejected = enum.auto()
    Unreachable = enum.auto()


class FlowStatus(enum.Enum):
    OK = enum.auto()
    Error = enum.auto()


TFlow = TypeVar("TFlow", bound="Flow")


class Flow(AbstractContextManager):
    def __init__(
        self,
        fcs_stub: FlowControlServiceStub,
        control_point: str,
        span: trace.Span,
        check_response: Optional[check_pb2.CheckResponse],
        ramp_mode: bool,
        cache_key: Optional[str],
    ):
        self._fcs_stub = fcs_stub
        self._control_point = control_point
        self._span = span
        self._check_response = check_response
        self._cache_key = cache_key
        self._status_code = FlowStatus.OK
        self._ended = False
        self._ramp_mode = ramp_mode
        self.logger = logging.getLogger("aperture-py-sdk-flow")

    def should_run(self) -> bool:
        return self.decision == FlowDecision.Accepted or (
            (not self._ramp_mode) and self.decision == FlowDecision.Unreachable
        )

    @property
    def decision(self) -> FlowDecision:
        if self.check_response is None:
            return FlowDecision.Unreachable
        if (
            self.check_response.decision_type
            == check_pb2.CheckResponse.DECISION_TYPE_ACCEPTED
        ):
            return FlowDecision.Accepted
        return FlowDecision.Rejected

    @property
    def success(self) -> bool:
        return self.decision != FlowDecision.Unreachable

    @property
    def check_response(self) -> Optional[check_pb2.CheckResponse]:
        return self._check_response

    def set_status(self, status_code: FlowStatus) -> None:
        self._status_code = status_code

    def end(self) -> None:
        if self._ended:
            self.logger.warning("attempting to end an already ended flow")
            return
        self._ended = True

        check_response_json = (
            json_format.MessageToJson(self.check_response)
            if self.check_response
            else ""
        )
        self._span.set_attributes(
            {
                flow_status_label: self._status_code.name,
                check_response_label: check_response_json,
                flow_end_timestamp_label: time.monotonic_ns(),
            }
        )
        self._span.end()

    def set_result_cache(self, value: str, ttl: datetime.timedelta):
        if not self._cache_key:
            return KeyUpsertResponse(ValueError("No cache key"))

        cache_upsert_request = {
            "controlPoint": self._control_point,
            "key": self._cache_key,
            "value": value,
            "ttl": ttl,
        }

        res: CacheUpsertResponse = self._fcs_stub.CacheUpsert(cache_upsert_request)

        if res.result_cache_response == None:
            return KeyUpsertResponse(ValueError("No cache upsert response"))

        return KeyUpsertResponse(
            convert_cache_error(res.result_cache_response.error),
        )

    async def delete_result_cache(self):
        if not self._cache_key:
            raise ValueError("No cache key")

        cache_delete_request = {
            "controlPoint": self._control_point,
            "key": self._cache_key,
        }

        res: CacheDeleteResponse = self._fcs_stub.CacheDelete(cache_delete_request)

        if res.result_cache_response == None:
            return KeyDeleteResponse(ValueError("No cache delete response"))

        return KeyDeleteResponse(
            convert_cache_error(res.result_cache_response.error),
        )

    def result_cache(self):
        if (
            not self.check_response
            or not self.check_response.cache_lookup_response
            or not self.check_response.cache_lookup_response.result_cache_response
        ):
            return KeyLookupResponse(None, MISS, ValueError("No cache lookup response"))

        lookup_response = (
            self.check_response.cache_lookup_response.result_cache_response
        )
        return KeyLookupResponse(
            lookup_response.value,
            convert_cache_lookup_status(lookup_response.lookup_status),
            None,
        )

    def __enter__(self: TFlow) -> TFlow:
        return self

    def __exit__(self, exc_type, _exc_value, _traceback) -> None:
        if self._ended:
            return
        if exc_type is not None:
            self.set_status(FlowStatus.Error)
        self.end()
