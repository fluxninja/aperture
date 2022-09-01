import http from "k6/http";
import { check } from "k6";

export let vuStages = [
  { duration: "30s", target: 5 }, // simulate ramp-up of traffic from 0 to 5 users over 30 seconds
  { duration: "30s", target: 5 }, // stay at 5 users for 30s minutes
  { duration: "2m", target: 15 }, // ramp-up to 10 users over 1 minutes
  { duration: "2m", target: 15 }, // stay at 10 users for 2 minutes (peak hour)
  { duration: "10s", target: 5 }, // ramp-down to 5 users in 10 seconds
  { duration: "30s", target: 5 }, // stay at to 5 users in 30 seconds
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
  const url = "http://demo1-demo-app.demoapp.svc.cluster.local/request";
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
