import { check, group, sleep } from 'k6';
import http from 'k6/http';

const BASE_URL = 'http://10.204.0.36:61619';  // Assuming ActiveMQ Artemis is running on port 8161
const QUEUE_NAME = 'TESTKUBE';  // Replace with the actual queue name
const ARTEMIS_CREDENTIALS = 'cgi:cgi';  // Replace with your credentials

export const options = {
  stages: [
    { duration: '1m', target: 50 }, // Ramp up to 50 virtual users over 1 minute
    { duration: '1m', target: 50 }, // Stay at 50 virtual users for 1 minute
    { duration: '1m', target: 0 },  // Ramp down to 0 virtual users over 1 minute
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests must complete within 500ms
  },
};

export default function () {
  group('Send Message to Queue', () => {
    // Prepare the message payload
    const messagePayload = 'hi, this is a test';

    console.log('Sending message to the queue:', QUEUE_NAME, 'Message:', messagePayload);

    // Send a message to the queue using ActiveMQ Artemis REST API
    const sendMessageResponse = http.post(
      `${BASE_URL}/api/message/${QUEUE_NAME}?type=queue`,
      `body=${messagePayload}`,
      {
        headers: {
          Authorization: `Basic ${ARTEMIS_CREDENTIALS}`,
        },
      }
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
