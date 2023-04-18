import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { check, sleep } from "k6";
import encoding from "k6/encoding";
import { vu } from "k6/execution";
import http from "k6/http";

export let vuStages = [
  { target: 600, duration: "1m" },
  { target: 5000, duration: "2m" },
  { target: 10000, duration: "4m" },
  { target: 600, duration: "5m" },
];

export let options = {
  discardResponseBodies: true,
  scenarios: {
    user: {
      executor: "ramping-arrival-rate",
      startRate: 2,
      timeUnit: "1m",
      preAllocatedVUs: 2,
      maxVUs: 50,
      stages: vuStages,
      env: { USER_TYPE: "user" },
    },
  },
};

export default function () {
  let userType = __ENV.USER_TYPE;
  let userId = vu.idInTest;
  const url = "http://service1-demo-app.demoapp.svc.cluster.local/request";
  let cookieJsonStr = JSON.stringify({ user_type: userType });
  let cookieSession = encoding.b64encode(cookieJsonStr);
  const headers = {
    "Content-Type": "application/json",
    Cookie: "session=".concat(cookieSession),
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
