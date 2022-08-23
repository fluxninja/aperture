import http from "k6/http";
import { check } from "k6";

export let options = {
  stages: [
    { duration: "2m", target: 100 }, // below normal load
    { duration: "5m", target: 100 },
    { duration: "2m", target: 200 }, // normal load
    { duration: "5m", target: 200 },
    { duration: "2m", target: 300 }, // around the breaking point
    { duration: "5m", target: 300 },
    { duration: "2m", target: 400 }, // beyond the breaking point
    { duration: "5m", target: 400 },
    { duration: "10m", target: 0 }, // scale down. Recovery stage.
  ],
};

export default function () {
  const url = "http://demo1-demo-app.demoapp.svc.cluster.local/request";
  const headers = {
    "Content-Type": "application/json",
    Cookie:
      "session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo",
  };
  const body = {
    request: [
      [
        { destination: "demo1-demo-app.demoapp.svc.cluster.local" },
        { destination: "demo2-demo-app.demoapp.svc.cluster.local" },
        { destination: "demo3-demo-app.demoapp.svc.cluster.local" },
        { destination: "demo2-demo-app.demoapp.svc.cluster.local" },
      ],
    ],
  };
  let res = http.request("POST", url, JSON.stringify(body), {
    headers: headers,
  });
  check(res, {
    "http status was 200": res.status === 200,
  });
}
