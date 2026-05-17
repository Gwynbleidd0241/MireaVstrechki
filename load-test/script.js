import http from "k6/http";
import { check, sleep } from "k6";
import { Rate } from "k6/metrics";

const errorRate = new Rate("errors");

export const options = {
  stages: [
    { duration: "30s", target: 20 },
    { duration: "1m", target: 50 },
    { duration: "30s", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"],
    errors: ["rate<0.01"],
  },
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";

export function setup() {
  const email = __ENV.LOGIN || "admin@demo.local";
  const password = __ENV.PASSWORD || "password";

  const res = http.post(
    `${BASE_URL}/login`,
    JSON.stringify({ email, password }),
    { headers: { "Content-Type": "application/json" } },
  );

  console.log(`login status: ${res.status}, body: ${res.body}`);
  check(res, { "login ok": (r) => r.status === 200 });
  return { token: res.json("token") };
}

export default function (data) {
  const headers = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${data.token}`,
  };

  const listRes = http.get(`${BASE_URL}/events`, { headers });
  const listOk = check(listRes, {
    "events list 200": (r) => r.status === 200,
  });
  errorRate.add(!listOk);

  sleep(0.5);

  const events = listRes.json();
  if (Array.isArray(events) && events.length > 0) {
    const id = events[0].id;

    const eventRes = http.get(`${BASE_URL}/events/${id}`, { headers });
    const eventOk = check(eventRes, {
      "event get 200": (r) => r.status === 200,
    });
    errorRate.add(!eventOk);

    const tasksRes = http.get(`${BASE_URL}/events/${id}/tasks`, { headers });
    const tasksOk = check(tasksRes, { "tasks 200": (r) => r.status === 200 });
    errorRate.add(!tasksOk);
  }

  sleep(1);
}
