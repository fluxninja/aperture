import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);

export const PROTO_PATH = path.resolve(
  path.dirname(__filename),
  "../proto/check/v1/check.proto",
);

const fn_host = process.env.FN_AGENT_HOST
  ? process.env.FN_AGENT_HOST
  : "localhost";
const fn_port = process.env.FN_AGENT_PORT ? process.env.FN_AGENT_PORT : "8089";
export const URL = fn_host + ":" + fn_port;

export const LIBRARY_NAME = "aperture-js";
export const LIBRARY_VERSION = "v0.1.0";

// Label to hold source of flow.
export const SOURCE_LABEL = "aperture.source";
// Label to hold status of the flow.
export const FLOW_STATUS_LABEL = "aperture.flow.status";
// Label to hold JSON encoded check response struct.
export const CHECK_RESPONSE_LABEL = "aperture.check_response";
// Label to hold flow's start timestamp in Unix nanoseconds since Epoch.
export const FLOW_START_TIMESTAMP_LABEL = "aperture.flow_start_timestamp";
// Label to hold flow's stop timestamp in Unix nanoseconds since Epoch.
export const FLOW_END_TIMESTAMP_LABEL = "aperture.flow_end_timestamp";
// Label to hold workload start timestamp in Unix nanoseconds since Epoch.
export const WORKLOAD_START_TIMESTAMP_LABEL =
  "aperture.workload_start_timestamp";
