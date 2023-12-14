import { check, group } from 'k6';
import http from 'k6/http';

const BASE_URL = 'http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619';

export const options = {
  stages: [
    { duration: '1m', target: 50 }, // Ramp up to 50 virtual users over 2 minutes
    { duration: '1m', target: 50 }, // Stay at 50 virtual users for 5 minutes
    { duration: '1m', target: 0 },  // Ramp down to 0 virtual users over 2 minutes
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests must complete within 500ms
  },
};

const credentials = { username: 'cgi', password: 'cgi' };
const queueName = 'TESTKUBE'; // Replace with the actual queue name

export default function () {
  group('Send Message to Queue', () => {
    // Prepare the message payload
    const messagePayload = JSON.stringify({ key: 'value' });

    console.log('Sending message to the queue:', queueName);

    // Send a message to the queue
    const sendMessageResponse = http.post(
      `${BASE_URL}/queues/${queueName}/send`,
      messagePayload,
      { auth: credentials }
    );

    console.log('Response status:', sendMessageResponse.status);

    // Check if the request was successful
    check(sendMessageResponse, {
      'Message Sent Successfully': (resp) => resp.status === 200,
    });

    // Sleep for a short duration to simulate some processing time
    sleep(0.5);
  });
}
