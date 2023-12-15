import { group, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.1.0/index.js';

const BASE_URL = 'http://10.204.1.8:61616'; // Adjust the protocol and IP as needed
const credentials = { username: 'artemis', password: 'artemis' };
const queueName = 'TESTKUBE'; // Replace with the actual queue name

export const options = {
  stages: [
    { duration: '1m', target: 50 },
    { duration: '1m', target: 50 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
  },
};

export default function () {
  group('Send Message to Queue', () => {
    // Prepare the message payload
    const messagePayload = `Hi, this is a test - ${randomString(10)}`;

    console.log('Sending message to the queue:', queueName, 'Message:', messagePayload);

    // Simulate sending a message to the queue using HTTP
    const sendMessageResponse = http.post(
      `${BASE_URL}/jolokia/exec/org.apache.activemq.artemis:broker=\"artemis-broker\",component=addresses,address=\"TESTKUBE\"/sendMessage(java.lang.String)`,
      { timeout: '60s', auth: credentials }
    );

    console.log('Response status:', sendMessageResponse.status);

    // Sleep for a short duration to simulate some processing time
    sleep(0.5);
  });
}
