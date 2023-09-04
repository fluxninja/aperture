import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import { check, sleep } from "k6";
import { vu } from "k6/execution";
import http from "k6/http";

export let vuStages = [
  { duration: "10s", target: 5 },
  { duration: "2m", target: 5 },
  { duration: "1m", target: 50 },
  { duration: "2m", target: 50 },
  { duration: "10s", target: 5 },
  { duration: "2m", target: 5 },
];

export let options = {
  discardResponseBodies: true,
  scenarios: {
    guests: {
      executor: "ramping-vus",
      stages: vuStages,
      env: { USER_TYPE: "guest" },
    },
    subscribers: {
      executor: "ramping-vus",
      stages: vuStages,
      env: { USER_TYPE: "subscriber" },
    },
    bots: {
      executor: "ramping-vus",
      stages: vuStages,
      env: { USER_TYPE: "bot" },
    },
  },
};

export default function () {
  let userType = __ENV.USER_TYPE;
  let userId = vu.idInTest;
  const url = "http://service1-demo-app.demoapp.svc.cluster.local/prometheus";
  const headers = {
    "Content-Type": "application/json",
    "User-Type": userType,
    "User-Id": userId,
  };
  let res = http.request("GET", url, JSON.stringify({}), {
    headers: headers,
  });
  const ret = check(res, {
    "http status was 200": res.status === 200,
  });
  if (!ret) {
    // sleep for 10ms to 25ms
    sleep(randomIntBetween(0.01, 0.025));
  }
}
