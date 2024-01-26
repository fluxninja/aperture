import enum
from typing import Optional

from aperture_sdk._gen.aperture.flowcontrol.check.v1.check_pb2 import FlowEndResponse


class FlowDecision(enum.Enum):
    Accepted = enum.auto()
    Rejected = enum.auto()
    Unreachable = enum.auto()


class FlowStatus(enum.Enum):
    OK = enum.auto()
    Error = enum.auto()


class EndResponse:
    def __init__(
        self, error: Optional[Exception], flow_end_response: Optional[FlowEndResponse]
    ):
        self.error = error
        self.flow_end_response = flow_end_response

    def get_error(self) -> Optional[Exception]:
        return self.error

    def get_flow_end_response(self) -> Optional[FlowEndResponse]:
        return self.flow_end_response
