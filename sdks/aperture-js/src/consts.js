import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);

export const PROTO_PATH = path.resolve(
  path.dirname(__filename),
  "./proto/flowcontrol/v1/flowcontrol.proto"
);

export const host = 'localhost';
export const port = process.env.FN_APP_PORT ? process.env.FN_APP_PORT : "8000";

const fn_host = process.env.FN_AGENT_HOST ? process.env.FN_AGENT_HOST : "localhost";
const fn_port = process.env.FN_AGENT_PORT ? process.env.FN_AGENT_PORT : "8089";
export const URL = fn_host + ":" + fn_port;
