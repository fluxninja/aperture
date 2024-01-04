import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import { check, sleep } from "k6";
import { vu } from "k6/execution";
import http from "k6/http";

export let vuStages = [
  { duration: "25s", target: 500 },
  { duration: "20m", target: 1000 },
];

export let options = {
  discardResponseBodies: true,
  scenarios: {
    guests_api_key1: {
      executor: "ramping-vus",
      stages: vuStages,
      env: {
        USER_TYPE: "guest",
        API_KEY: "key1",
      },
    },
    subscribers_api_key1: {
      executor: "ramping-vus",
      stages: vuStages,
      env: {
        USER_TYPE: "subscriber",
        API_KEY: "key1",
      },
    },
    guests_api_key2: {
      executor: "ramping-vus",
      stages: vuStages,
      env: {
        USER_TYPE: "guest",
        API_KEY: "key2",
      },
    },
    subscribers_api_key2: {
      executor: "ramping-vus",
      stages: vuStages,
      env: {
        USER_TYPE: "subscriber",
        API_KEY: "key2",
      },
    },
  },
};

export default function () {
  let userType = __ENV.USER_TYPE;
  let userId = vu.idInTest;
  let apiKey = __ENV.API_KEY;
  const url = "http://service1-demo-app.demoapp.svc.cluster.local/request";
  const headers = {
    "Content-Type": "application/json",
    Cookie:
      "session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo",
    "User-Type": userType,
    "User-Id": userId,
    "Api-Key": apiKey,
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
