import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { check, sleep } from "k6";
import { vu } from "k6/execution";
import http from "k6/http";

export let vuStagesAPI = [
  { duration: "10s", target: 5 },
  { duration: "2m", target: 5 },
  { duration: "1m", target: 50 },
  { duration: "2m", target: 50 },
  { duration: "10s", target: 5 },
  { duration: "2m", target: 5 },
];

export let options = {
  insecureSkipTLSVerify: true,
  discardResponseBodies: true,
  scenarios: {
    api: {
      executor: "ramping-vus",
      stages: vuStagesAPI,
      env: { USER_TYPE: "api" },
    },
    agent: {
      executor: "ramping-vus",
      stages: vuStagesAPI,
      env: { USER_TYPE: "agent" },
    },
  },
};

export let token = "";

export function refresh_token() {
  console.log("refreshing token");
  const url = "http://authn.cloud/login/basic-jwt?orgId=dcde9fb7-6cec-4a4a-b015-114795a65ed0";
  const body = "email=ann.place@placeholder.com&password=ann.place";

  let res = http.request("POST", url, body, {
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
  });

  let auth = res.headers["Authorization"]

  if (auth) {
    token = auth.split(" ")[1];
  }
};

export function agent_service_request() {
  const agent_service_url = "https://agent-service.cloud:80/fluxninja/v1/report";
  const agent_service_headers = {
    "Content-Type": "application/json",
    "apiKey": "e85db05033e848d6815fa916651323c1",
  };
  const agent_service_body = {
    "version_info": {
      "version": "1.0.0",
      "service": "aperture-controller",
      "build_host": "",
      "build_os": "",
      "build_time": "",
      "git_branch": "",
      "git_commit_hash": "",
    },
    "host_info": {
      "hostname": "k6",
      "local_ip": "127.0.0.1",
      "uuid": "k6",
    },
    "agent_group": "default",
    "controller_info": {
      "id": "test_controller_" + vu.idInTest,
    },
    "installation_mode": "LINUX_BARE_METAL",
  };

  let res = http.request("POST", agent_service_url, JSON.stringify(agent_service_body), {
    headers: agent_service_headers,
  });

  const ret = check(res, {
    "http status was 200": res.status === 200,
  });
  if (!ret) {
    // sleep for 10ms to 25ms
    sleep(randomIntBetween(0.01, 0.025));
  }
};

export function api_service_request() {
  if (token == "") {
    refresh_token();
  }

  const api_service_url = "http://api-service.cloud:80/api/query";
  const api_service_headers = {
    "Content-Type": "application/json",
    "Authorization": "Bearer " + token,
  };
  const api_service_body = {
    "query": "\\n  query PoliciesGroupData(\\n    $first: Int\\n    $after: String\\n    $before: String\\n    $where: PolicyBoolExp\\n    $orderBy: [PolicyOrderBy]\\n  ) {\\n    policies(\\n      first: $first\\n      after: $after\\n      before: $before\\n      where: $where\\n      orderBy: $orderBy\\n    ) {\\n      totalCount\\n      pageInfo {\\n        endCursor\\n      }\\n      edges {\\n        node {\\n          id\\n          name\\n          circuit {\\n            node {\\n              componentId\\n            }\\n          }\\n          body\\n          end\\n          hash\\n          start\\n          status\\n          controller {\\n            name\\n            id\\n          }\\n        }\\n      }\\n    }\\n  }\\n",
    "variables": {
      "where": {
        "projectID": {
          "eq": "29d91ad4-9e46-404f-938a-7884c87c7523",
        },
        "not": {
          "status": {
            "eq": "3-archived",
          },
        },
      },
      "orderBy": {
        "name": "asc"
      },
      "first": 10,
    },
    "operationName": "PoliciesGroupData",
  };

  let res = http.request("POST", api_service_url, JSON.stringify(api_service_body), {
    headers: api_service_headers,
  });

  if (res.status == 401) {
    refresh_token();
    api_service_request();
  }

  const ret = check(res, {
    "http status was 200": res.status === 200,
  });
  if (!ret) {
    // sleep for 10ms to 25ms
    sleep(randomIntBetween(0.01, 0.025));
  }
};

export function setup() {
  refresh_token();
}

export default function () {
  let userType = __ENV.USER_TYPE;

  if (userType == "agent") {
    agent_service_request();
  } else {
    api_service_request();
  }
}
