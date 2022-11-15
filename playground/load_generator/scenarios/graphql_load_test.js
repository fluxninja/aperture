import http from "k6/http";
import { check, sleep } from "k6";
import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { vu } from "k6/execution";

export let options = {
    discardResponseBodies: true,
    scenarios: {
        contacts: {
            executor: 'constant-arrival-rate',
            duration: '2m',
            rate: 50,
            timeUnit: '1s',
            preAllocatedVUs: 1,
            maxVUs: 1,
        },
    },
};

export default function () {
    let userId = vu.idInTest;
    const url = "http://service-graphql-demo-app.demoapp.svc.cluster.local/query";
    const headers = {
        "Content-Type": "application/json",
        Cookie: "session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo",
    };
    const mutation = `mutation createTodo {
        createTodo(input: { text: "todo", userId: "${userId}" }) {
          user {
            id
          }
          text
          done
        }
      }`;
    let res = http.request("POST", url, JSON.stringify({ query: mutation }), {
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
