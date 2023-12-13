import http from 'k6/http';
import { check } from 'k6';

export let options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  // Replace the URL, username, and password with your Artemis broker endpoint and credentials
  let url = 'http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619';
  let username = 'cgi';
  let password = 'cgi';

  // Example: Sending a message to a queue
  let payload = JSON.stringify({
    message: 'Hello, Artemis!',
  });

  let params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Basic ${btoa(`${username}:${password}`)}`, // Basic Authentication
    },
  };

  let res = http.post(`${url}/queue/your-queue-name`, payload, params);

  check(res, {
    'message sent successfully': (r) => r.status === 200,
  });
}
