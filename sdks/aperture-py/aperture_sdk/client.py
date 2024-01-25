"""ApertureClient for starting Flows."""

import datetime
import functools
import logging
import time
import typing
from typing import Callable, Optional, Type, TypeVar

import grpc
from aperture_sdk._gen.aperture.flowcontrol.check.v1.check_pb2 import (
    CacheLookupRequest,
    CheckRequest,
)
from aperture_sdk._gen.aperture.flowcontrol.check.v1.check_pb2_grpc import (
    FlowControlServiceStub,
)
from aperture_sdk.client_common import *
from aperture_sdk.const import (
    default_grpc_reconnection_time,
    flow_start_timestamp_label,
    library_name,
    library_version,
    source_label,
    workload_start_timestamp_label,
)
from aperture_sdk.flow import Flow
from aperture_sdk.utils import TWrappedReturn, run_fn
from opentelemetry import baggage, trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import SERVICE_NAME, SERVICE_VERSION, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.util import types as otel_types

TApertureClient = TypeVar("TApertureClient", bound="ApertureClient")
TWrappedFunction = Callable[..., TWrappedReturn]


class ApertureClient:
    """
    ApertureClient can be used to start Flows.
    """

    def __init__(
        self,
        channel: grpc.Channel,
        otlp_exporter: OTLPSpanExporter,
    ):
        self.logger = logging.getLogger("aperture-py")

        resource = Resource.create(
            {
                SERVICE_NAME: library_name,
                SERVICE_VERSION: library_version,
            }
        )
        tracer_provider = TracerProvider(resource=resource)
        trace.set_tracer_provider(tracer_provider)
        self.tracer = trace.get_tracer(library_name, library_version)

        span_processor = BatchSpanProcessor(otlp_exporter)
        tracer_provider.add_span_processor(span_processor)
        self.otlp_exporter = otlp_exporter
        self.grpc_channel = channel

    @classmethod
    def new_client(
        cls: Type[TApertureClient],
        address: str,
        api_key: Optional[str] = None,
        insecure: bool = False,
        grpc_timeout: datetime.timedelta = default_grpc_reconnection_time,
        credentials: Optional[grpc.ChannelCredentials] = None,
        compression: grpc.Compression = grpc.Compression.NoCompression,
    ) -> TApertureClient:
        if not address:
            raise ValueError("Address must be provided")
        if not credentials:
            credentials = grpc.ssl_channel_credentials()
        if api_key:
            metadata_plugin_instance = ApertureCloudAuthMetadataPlugin(api_key)
            credentials = grpc.composite_channel_credentials(
                credentials,
                grpc.metadata_call_credentials(
                    metadata_plugin=metadata_plugin_instance,
                    name="x-api-key",
                ),
            )

        otlp_exporter = OTLPSpanExporter(
            endpoint=address,
            insecure=insecure,
            credentials=credentials,
            compression=compression,
            timeout=int(grpc_timeout.total_seconds()),
        )
        grpc_channel_options_dict = {
            "grpc.keepalive_time_ms": 10000,
            "grpc.keepalive_timeout_ms": 5000,
        }
        grpc_channel_options = [(k, v) for k, v in grpc_channel_options_dict.items()]
        grpc_channel = (
            grpc.insecure_channel(
                address, compression=compression, options=grpc_channel_options
            )
            if insecure
            else grpc.secure_channel(
                address,
                credentials,
                compression=compression,
                options=grpc_channel_options,
            )
        )
        return cls(
            channel=grpc_channel,
            otlp_exporter=otlp_exporter,
        )

    def start_flow(
        self,
        control_point: str,
        params: FlowParams,
    ) -> Flow:
        labels: Labels = {}
        labels.update({key: str(value) for key, value in baggage.get_all().items()})
        # Explicit labels override baggage
        labels.update(params.explicit_labels or {})
        request = CheckRequest(
            control_point=control_point,
            labels=labels,
            ramp_mode=params.ramp_mode,
            expect_end=True,
            cache_lookup_request=CacheLookupRequest(
                result_cache_key=params.result_cache_key,
                global_cache_keys=params.global_cache_keys,
            ),
        )
        span_attributes: otel_types.Attributes = {
            flow_start_timestamp_label: time.monotonic_ns(),
            source_label: "sdk",
        }

        span = self.tracer.start_span("Aperture Check", attributes=span_attributes)
        stub = FlowControlServiceStub(self.grpc_channel)
        error: Optional[Exception] = None
        try:
            # stub.Check is typed to accept an int, but it actually accepts a float
            timeout = typing.cast(int, params.check_timeout.total_seconds())
            if timeout == 0:
                timeout = None
            response = stub.Check(request, timeout=timeout)
        except grpc.RpcError as e:
            self.logger.debug(f"Aperture gRPC call failed: {e.details()}")
            response = None
            error = e
        span.set_attribute(workload_start_timestamp_label, time.monotonic_ns())
        return Flow(
            fcs_stub=stub,
            control_point=control_point,
            span=span,
            check_response=response,
            ramp_mode=params.ramp_mode,
            cache_key=params.result_cache_key,
            error=error,
            grpc_channel=self.grpc_channel,
        )

    def decorate(
        self,
        control_point: str,
        params: FlowParams = FlowParams(),
        on_reject: Optional[Callable] = None,
    ) -> Callable[[TWrappedFunction], TWrappedFunction]:
        def decorator(fn: TWrappedFunction) -> TWrappedFunction:
            @functools.wraps(fn)
            async def wrapper(*args, **kwargs):
                flow = self.start_flow(control_point, params)
                if flow.should_run():
                    return await run_fn(fn, *args, **kwargs)
                else:
                    if on_reject:
                        return on_reject()

            return wrapper

        return decorator

    def close(self):
        self.otlp_exporter.shutdown()
