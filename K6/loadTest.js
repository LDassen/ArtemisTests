import { check, sleep } from 'k6';
import http from 'k6/http';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  const url = 'http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619';  // Adjust the URL according to your Artemis setup
  const queueName = 'exampleQueueCore';
  const message = 'Hello, Core!';
  const username = 'cgi';
  const password = 'cgi';

  const payload = JSON.stringify({
    queueName,
    message,
    user: {
      username,
      password,
    },
  });

  const headers = {
    'Content-Type': 'application/json',
  };

  const response = http.post(`${url}/send-receive-endpoint`, payload, { headers });
  console.log(response.status, response.body);

  check(response, {
    'HTTP Request Successful': (r) => r.status === 200,
    'Core Message Received Successfully': (r) => r.json('receivedMessage') !== null,
  });

  // Sleep for a short duration between requests (adjust as needed)
  sleep(1);
}
