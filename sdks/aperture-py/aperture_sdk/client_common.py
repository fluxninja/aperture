import datetime
from dataclasses import dataclass
from typing import Dict, List, Optional

from aperture_sdk.const import default_rpc_timeout
from grpc import AuthMetadataContext, AuthMetadataPlugin, AuthMetadataPluginCallback

Labels = Dict[str, str]


class ApertureCloudAuthMetadataPlugin(AuthMetadataPlugin):
    def __init__(self, api_key):
        self.api_key = api_key

    def __call__(
        self, context: AuthMetadataContext, callback: AuthMetadataPluginCallback
    ) -> None:
        callback((("x-api-key", self.api_key),), None)


@dataclass
class FlowParams:
    explicit_labels: Optional[Labels] = None
    check_timeout: datetime.timedelta = default_rpc_timeout
    ramp_mode: bool = False
    result_cache_key: Optional[str] = None
    global_cache_keys: Optional[List[str]] = None
