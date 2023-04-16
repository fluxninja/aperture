import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { check, sleep } from "k6";
import { vu } from "k6/execution";
import http from "k6/http";

export let vuStages = [
  { duration: "10s", target: 5 },
  { duration: "2m", target: 5 },
  { duration: "1m", target: 30 },
  { duration: "2m", target: 30 },
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
  const url = "http://service1-demo-app.demoapp.svc.cluster.local/request";
  const headers = {
    "Content-Type": "application/json",
    Cookie:
      "session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo",
    "User-Type": userType,
    "User-Id": userId,
  };
  const body = {
    request: [
      [
        {
          destination: "service1-demo-app.demoapp.svc.cluster.local/request",
        },
        {
          destination: "service2-demo-app.demoapp.svc.cluster.local/request",
        },
        {
          destination: "service3-demo-app.demoapp.svc.cluster.local/request",
        },
      ],
    ],
  };
  let res = http.request("POST", url, JSON.stringify(body), {
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
