import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { check, sleep } from "k6";
import { vu } from "k6/execution";
import http from "k6/http";

// export let vuStages = [
//   { duration: "10s", target: 50 },
//   { duration: "2m", target: 50 },
//   { duration: "1m", target: 50 },
//   { duration: "2m", target: 50 },
//   { duration: "10s", target: 5 },
//   { duration: "2m", target: 5 },
// ];

export let options = {
  discardResponseBodies: true,
  scenarios: {
    guests: {
      executor: "constant-arrival-rate",
      duration: '5m',
      rate: 2000,
      timeUnit: '1s',
      preAllocatedVUs: 100,
      env: { USER_TYPE: "guest" },
    },
    subscribers: {
      executor: "constant-arrival-rate",
      duration: '5m',
      rate: 1000,
      timeUnit: '1s',
      preAllocatedVUs: 100,
      env: { USER_TYPE: "subscriber" },
    },
    crawlers: {
      executor: "constant-arrival-rate",
      duration: '5m',
      rate: 2000,
      timeUnit: '1s',
      preAllocatedVUs: 100,
      env: { USER_TYPE: "crawler" },
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
    "One-Header": "one",
    "Two-Header": "two",
    "Three-Header": "three",
    "Four-Header": "four",
    "Five-Header": "five",
    "Six-Header": "six",
    "Seven-Header": "seven",
    "Eight-Header": "eight",
    "Nine-Header": "nine",
    "Ten-Header": "ten",
    "One-One-Header": "one-one",
    "One-Two-Header": "one-two",
    "One-Three-Header": "one-three",
    "One-Four-Header": "one-four",
    "One-Five-Header": "one-five",
    "One-Six-Header": "one-six",
    "One-Seven-Header": "one-seven",
    "One-Eight-Header": "one-eight",
    "One-Nine-Header": "one-nine",
    "One-Ten-Header": "one-ten",
    "Two-One-Header": "two-one",
    "Two-Two-Header": "two-two",
    "Two-Three-Header": "two-three",
    "Two-Four-Header": "two-four",
    "Two-Five-Header": "two-five",
    "Two-Six-Header": "two-six",
    "Two-Seven-Header": "two-seven",
    "Two-Eight-Header": "two-eight",
    "Two-Nine-Header": "two-nine",
    "Two-Ten-Header": "two-ten",
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
