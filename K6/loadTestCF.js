import { check, group, sleep } from 'k6';
import { ConnectionFactory, TextMessage } from 'k6/x/jms';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.1.0/index.js';

const BASE_URL = 'tcp://10.204.1.8:61616'; // Adjust the protocol and IP as needed
const credentials = { username: 'artemis', password: 'artemis' };
const queueName = 'TESTKUBE'; // Replace with the actual queue name

const connectionFactory = new ConnectionFactory({
  brokerURL: BASE_URL,
  credentials: credentials,
});

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
    const messagePayload = `Hi, this is a test - ${randomString(10)}`;

    console.log('Sending message to the queue:', queueName, 'Message:', messagePayload);

    // Send a message to the queue using JMS
    const connection = connectionFactory.createConnection();
    connection.start();

    try {
      const session = connection.createSession();
      const producer = session.createProducer(session.createQueue(queueName));

      const message = new TextMessage(messagePayload);
      producer.send(message);

      console.log('Message sent successfully.');
    } finally {
      connection.close();
    }

    // Sleep for a short duration to simulate some processing time
    sleep(0.5);
  });
}
