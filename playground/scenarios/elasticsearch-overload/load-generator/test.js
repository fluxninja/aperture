import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import { check, sleep } from "k6";
import encoding from "k6/encoding";
import { vu } from "k6/execution";
import http from "k6/http";

export let vuStages = [
    { duration: "0s", target: 50 },
    { duration: "1m", target: 50 },
    { duration: "1m", target: 250 },
    { duration: "10m", target: 250 },
    { duration: "1m", target: 50 },
    { duration: "1m", target: 0 },
];

export let options = {
    discardResponseBodies: true,
    scenarios: {
        subscribers: {
            executor: "ramping-vus",
            stages: vuStages,
            env: { USER_TYPE: "subscriber" },
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
    const body = {};
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
