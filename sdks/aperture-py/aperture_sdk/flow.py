"""Flow started using the ApertureClient."""

import datetime
import logging
import time
from contextlib import AbstractContextManager
from typing import List, Optional, TypeVar

import grpc
from aperture_sdk._gen.aperture.flowcontrol.check.v1 import check_pb2
from aperture_sdk._gen.aperture.flowcontrol.check.v1.check_pb2 import (
    CacheDeleteRequest,
    CacheDeleteResponse,
    CacheEntry,
    CacheUpsertRequest,
    FlowEndRequest,
    InflightRequestRef,
    StatusCode,
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
from aperture_sdk.flow_common import *
from google.protobuf import json_format
from google.protobuf.duration_pb2 import Duration
from opentelemetry import trace

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
        error: Optional[Exception],
        grpc_channel: grpc.Channel,
    ):
        self._fcs_stub = fcs_stub
        self._control_point = control_point
        self._span = span
        self._check_response = check_response
        self._cache_key = cache_key
        self._status_code = FlowStatus.OK
        self._ended = False
        self._ramp_mode = ramp_mode
        self._error = error
        self._grpc_channel = grpc_channel
        self.logger = logging.getLogger("aperture-py-sdk-flow")

    def should_run(self) -> bool:
        return self.decision == FlowDecision.Accepted or (
            (not self._ramp_mode) and self.decision == FlowDecision.Unreachable
        )

    def http_response_code(self) -> Optional[int]:
        if self.check_response is None:
            return 200  # Default to 200 if check_response is None

        if self.check_response.denied_response_status_code == StatusCode.Empty:
            return 200  # Default to 200 if denied_response_status_code is None or Empty

        return self.check_response.denied_response_status_code

    def retry_after(self) -> Optional[datetime.timedelta]:
        if self.check_response is None:
            return None
        if self.check_response.wait_time is None:
            return None
        return datetime.timedelta(seconds=self.check_response.wait_time.seconds)

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

    def end(self) -> EndResponse:
        if self._ended:
            self.logger.warning("attempting to end an already ended flow")
            return EndResponse(
                error=Exception("attempting to end an already ended flow"),
                flow_end_response=None,
            )
        if not self.check_response:
            self.logger.warning("attempting to end a flow with no check response")
            return EndResponse(
                error=Exception("attempting to end a flow with no check response"),
                flow_end_response=None,
            )
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

        inflight_request_ref: List[InflightRequestRef] = list()

        if self.check_response:
            for decision in self.check_response.limiter_decisions:
                if decision.WhichOneof("details") == "concurrency_limiter_info":
                    if decision.concurrency_limiter_info.request_id == "":
                        continue
                    ref: InflightRequestRef = InflightRequestRef(
                        policy_name=decision.policy_name,
                        policy_hash=decision.policy_hash,
                        component_id=decision.component_id,
                        label=decision.concurrency_limiter_info.label,
                        request_id=decision.concurrency_limiter_info.request_id,
                    )
                    if decision.concurrency_limiter_info.tokens_info:
                        ref.tokens = (
                            decision.concurrency_limiter_info.tokens_info.consumed
                        )
                    inflight_request_ref.append(ref)
                elif decision.WhichOneof("details") == "concurrency_scheduler_info":
                    if decision.concurrency_scheduler_info.request_id == "":
                        continue
                    ref: InflightRequestRef = InflightRequestRef(
                        policy_name=decision.policy_name,
                        policy_hash=decision.policy_hash,
                        component_id=decision.component_id,
                        label=decision.concurrency_scheduler_info.label,
                        request_id=decision.concurrency_scheduler_info.request_id,
                    )
                    if decision.concurrency_scheduler_info.tokens_info:
                        ref.tokens = (
                            decision.concurrency_scheduler_info.tokens_info.consumed
                        )
                    inflight_request_ref.append(ref)

        if inflight_request_ref:
            flow_end_request = FlowEndRequest(
                control_point=self._control_point,
                inflight_requests=inflight_request_ref,
            )

            try:
                res = self._fcs_stub.FlowEnd(flow_end_request)
            except grpc.RpcError as e:
                self.logger.error(f"Aperture gRPC call failed: {e.details()}")
                return EndResponse(
                    error=e,
                    flow_end_response=None,
                )

            return EndResponse(
                error=None,
                flow_end_response=res,
            )
        else:
            return EndResponse(
                error=None,
                flow_end_response=None,
            )

    def error(self) -> Optional[Exception]:
        return self._error

    def set_result_cache(
        self, value: str, ttl: datetime.timedelta, **grpc_opts
    ) -> KeyUpsertResponse:
        if not self._cache_key:
            return KeyUpsertResponse(ValueError("No cache key"))

        ttl_duration = Duration()
        ttl_duration.FromTimedelta(ttl)
        cache_upsert_request = CacheUpsertRequest(
            control_point=self._control_point,
            result_cache_entry=CacheEntry(
                key=self._cache_key,
                value=bytes(value, "utf-8"),
                ttl=ttl_duration,
            ),
        )

        try:
            res = self._fcs_stub.CacheUpsert(cache_upsert_request, **grpc_opts)
        except grpc.RpcError as e:
            self.logger.debug(f"Aperture gRPC call failed: {e.details()}")
            return KeyUpsertResponse(e)

        if res.result_cache_response is None:
            return KeyUpsertResponse(ValueError("No cache upsert response"))

        return KeyUpsertResponse(
            convert_cache_error(res.result_cache_response.error),
        )

    def delete_result_cache(self, **grpc_opts) -> KeyDeleteResponse:
        if not self._cache_key:
            return KeyDeleteResponse(ValueError("No cache key"))

        cache_delete_request = CacheDeleteRequest(
            control_point=self._control_point,
            result_cache_key=self._cache_key,
        )

        try:
            res: CacheDeleteResponse = self._fcs_stub.CacheDelete(
                cache_delete_request, **grpc_opts
            )
        except grpc.RpcError as e:
            self.logger.debug(f"Aperture gRPC call failed: {e.details()}")
            return KeyDeleteResponse(e)

        if res.result_cache_response is None:
            return KeyDeleteResponse(ValueError("No cache delete response"))

        return KeyDeleteResponse(
            convert_cache_error(res.result_cache_response.error),
        )

    def result_cache(self) -> KeyLookupResponse:
        if self._error is not None:
            return KeyLookupResponse(None, LookupStatus.MISS, self._error)
        if not self.check_response:
            return KeyLookupResponse(
                None, LookupStatus.MISS, ValueError("response is null")
            )
        if not self.should_run():
            return KeyLookupResponse(
                None, LookupStatus.MISS, ValueError("flow was rejected")
            )

        if (
            not self.check_response.cache_lookup_response
            or not self.check_response.cache_lookup_response.result_cache_response
        ):
            return KeyLookupResponse(
                None, LookupStatus.MISS, ValueError("No cache lookup response")
            )
        lookup_response = (
            self.check_response.cache_lookup_response.result_cache_response
        )
        return KeyLookupResponse(
            lookup_response.value,
            convert_cache_lookup_status(lookup_response.lookup_status),
            convert_cache_error(lookup_response.error),
        )

    def set_global_cache(
        self, key: str, value: str, ttl: datetime.timedelta, **grpc_opts
    ) -> KeyUpsertResponse:
        ttl_duration = Duration()
        ttl_duration.FromTimedelta(ttl)
        cache_upsert_request = CacheUpsertRequest(
            global_cache_entries={
                key: CacheEntry(
                    value=bytes(value, "utf-8"),
                    ttl=ttl_duration,
                ),
            },
        )

        try:
            res = self._fcs_stub.CacheUpsert(cache_upsert_request, **grpc_opts)
        except grpc.RpcError as e:
            self.logger.debug(f"Aperture gRPC call failed: {e.details()}")
            return KeyUpsertResponse(e)

        responses = res.global_cache_responses
        if responses is None:
            return KeyUpsertResponse(ValueError("No cache upsert response"))
        if key not in responses:
            return KeyUpsertResponse(
                ValueError("Key missing from global cache response")
            )

        return KeyUpsertResponse(
            convert_cache_error(responses[key].error),
        )

    def delete_global_cache(self, key: str, **grpc_opts) -> KeyDeleteResponse:
        cache_delete_request = CacheDeleteRequest(
            global_cache_keys=[key],
        )

        try:
            res: CacheDeleteResponse = self._fcs_stub.CacheDelete(
                cache_delete_request, **grpc_opts
            )
        except grpc.RpcError as e:
            self.logger.debug(f"Aperture gRPC call failed: {e.details()}")
            return KeyDeleteResponse(e)

        delete_responses = res.global_cache_responses

        if delete_responses is None:
            return KeyDeleteResponse(ValueError("No cache delete response"))
        if key not in delete_responses:
            return KeyDeleteResponse(
                ValueError("Key missing from global cache response")
            )

        return KeyDeleteResponse(
            convert_cache_error(delete_responses[key].error),
        )

    def global_cache(self, key: str) -> KeyLookupResponse:
        if self._error is not None:
            return KeyLookupResponse(None, LookupStatus.MISS, self._error)
        if not self.check_response:
            return KeyLookupResponse(
                None, LookupStatus.MISS, ValueError("response is null")
            )
        if not self.should_run():
            return KeyLookupResponse(
                None, LookupStatus.MISS, ValueError("flow was rejected")
            )

        if (
            not self.check_response.cache_lookup_response
            or not self.check_response.cache_lookup_response.global_cache_responses
        ):
            return KeyLookupResponse(
                None, LookupStatus.MISS, ValueError("No global cache lookup response")
            )

        lookup_response_map = (
            self.check_response.cache_lookup_response.global_cache_responses
        )
        if key not in lookup_response_map:
            return KeyLookupResponse(
                None, LookupStatus.MISS, ValueError("Unknown global cache key")
            )

        lookup_response = lookup_response_map[key]
        return KeyLookupResponse(
            lookup_response.value,
            convert_cache_lookup_status(lookup_response.lookup_status),
            convert_cache_error(lookup_response.error),
        )

    def __enter__(self: TFlow) -> TFlow:
        return self

    def __exit__(self, exc_type, _exc_value, _traceback) -> None:
        if self._ended:
            return
        if exc_type is not None:
            self.set_status(FlowStatus.Error)
        res = self.end()

        if res.get_error():
            self.logger.warning(f"Failed to end flow: {res.get_error()}")

        if res.get_flow_end_response():
            print(f"Ended flow: {res.get_flow_end_response()}")
