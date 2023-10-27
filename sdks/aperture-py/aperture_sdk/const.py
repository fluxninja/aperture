from datetime import timedelta

# Library name and version can be used by the user to create a resource that connects to telemetry exporter.
library_name = "aperture-py"
library_version = "v2.22.0"

# Config defaults.
default_rpc_timeout = timedelta(milliseconds=200)
default_grpc_reconnection_time = timedelta(seconds=10)

# Label keys.
# Label to hold source of flow.
source_label = "aperture.source"
# Label to hold status of the flow.
flow_status_label = "aperture.flow.status"
# Label to hold JSON encoded check response struct.
check_response_label = "aperture.check_response"
# Label to hold flow's start timestamp in Unix nanoseconds since Epoch.
flow_start_timestamp_label = "aperture.flow_start_timestamp"
# Label to hold flow's stop timestamp in Unix nanoseconds since Epoch.
flow_end_timestamp_label = "aperture.flow_end_timestamp"
# Label to hold workload start timestamp in Unix nanoseconds since Epoch.
workload_start_timestamp_label = "aperture.workload_start_timestamp"
