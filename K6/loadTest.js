import http from 'k6/http';
import { check } from 'k6';

export let options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  // Replace the URL with your Artemis broker endpoint
  let url = 'http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619';

  // Example: Sending a message to a queue
  let payload = JSON.stringify({
    message: 'Hello, Artemis!',
  });

  let params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  let res = http.post(`${url}/queue/your-queue-name`, payload, params);

  check(res, {
    'message sent successfully': (r) => r.status === 200,
  });
}
