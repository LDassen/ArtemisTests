import { check } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  const url = 'http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619';  // Adjust the URL according to your Artemis setup
  const queueName = 'exampleQueueCore';
  const message = 'Hello, Core!';

  const response = http.post(`${url}/send-receive-endpoint`, JSON.stringify({ queueName, message }), {
    headers: {
      'Content-Type': 'application/json',
    },
  });

  check(response, {
    'HTTP Request Successful': (r) => r.status === 200,
    'Core Message Received Successfully': (r) => r.json('receivedMessage') !== null,
  });
}
