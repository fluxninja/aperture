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
        user: {
            executor: "ramping-vus",
            stages: vuStages,
            env: { USER_TYPE: "user" },
        },
    },
};

export default function () {
    let userType = __ENV.USER_TYPE;
    let userId = vu.idInTest;

    const headers = {
        'Content-Type': 'application/json',
        'Cookie': 'session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo',
        'User-Type': userType,
        'User-Id': userId,
    };


    const svc1Body = {
        request: [
            [
                {
                    destination: "service1-demo-app.demoapp.svc.cluster.local/request",
                },
                {
                    destination: "service3-demo-app.demoapp.svc.cluster.local/request",
                },
            ],
        ],
    };
    const svc1 = {
        method: 'POST',
        url: 'http://service1-demo-app.demoapp.svc.cluster.local/request',
        body: JSON.stringify(svc1Body),
        params: {
            headers: headers,
        },
    };

    const svc2Body = {
        request: [
            [
                {
                    destination: "service2-demo-app.demoapp.svc.cluster.local/request",
                },
                {
                    destination: "service3-demo-app.demoapp.svc.cluster.local/request",
                },
            ],
        ],
    };
    const svc2 = {
        method: 'POST',
        url: 'http://service2-demo-app.demoapp.svc.cluster.local/request',
        body: JSON.stringify(svc2Body),
        params: {
            headers: headers,
        },
    };

    const responses = http.batch([svc1, svc2]);

    const svc1ret = check(responses[0], {
        "http status was 200": responses[0].status === 200,
    });
    const svc2ret = check(responses[1], {
        "http status was 200": responses[1].status === 200,
    });

    if (!svc1ret && !svc2ret) {
        // sleep for 10ms to 25ms
        sleep(randomIntBetween(10, 25) / 1000);
    };
}
