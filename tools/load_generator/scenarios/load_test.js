import http from "k6/http";
import { check } from "k6";

export let vuStages = [
  { duration: "1s", target: 5 },
  { duration: "2m", target: 5 },
  { duration: "1m", target: 30 },
  { duration: "2m", target: 30 },
  { duration: "1s", target: 5 },
  { duration: "5m", target: 5 },
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
  },
};

export default function () {
  let userType = __ENV.USER_TYPE;
  const url = "http://service1-demo-app.demoapp.svc.cluster.local/request";
  const headers = {
    "Content-Type": "application/json",
    Cookie:
      "session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo",
    "User-Type": userType,
  };
  const body = {};
  let res = http.request("POST", url, JSON.stringify(body), {
    headers: headers,
  });
  check(res, {
    "http status was 200": res.status === 200,
  });
}
