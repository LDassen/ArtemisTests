import http from 'k6/http';

const BASE_URL = 'http://10.204.0.47:61619';  // Correct Artemis management URL
const QUEUE_NAME = 'TESTKUBE';              // Replace with your queue name
const ARTEMIS_CREDENTIALS = 'cgi:cgi';      // Replace with your credentials

export default function () {
  // Example JSON payload
  const messagePayload = {
    text: 'hi, this is a test',
  };
  console.log('Message Payload:', JSON.stringify(messagePayload));

  // Send a message to the queue using Artemis management API
  const sendMessageResponse = http.post(
    `${BASE_URL}/console/jolokia/exec/org.apache.activemq.artemis:broker="0.0.0.0",component=addresses,address="${QUEUE_NAME}",subcomponent=queues,routing-type="anycast",queue="${QUEUE_NAME}"/sendMessage`,
    JSON.stringify(messagePayload),
    {
      headers: {
        'Content-Type': 'application/json',  // Set the content type for JSON payload
        Authorization: `Basic ${ARTEMIS_CREDENTIALS}`,
      },
    }
  );

  console.log('Request URL:', `${BASE_URL}/console/jolokia/exec/...`);
  console.log('Request Headers:', { Authorization: `Basic ${ARTEMIS_CREDENTIALS}` });
  console.log('Response status:', sendMessageResponse.status);
  console.log('Response body:', sendMessageResponse.body);

  // Check if the request was successful
  if (sendMessageResponse.status === 200) {
    console.log('Message sent successfully');
  } else {
    console.error('Failed to send message:', sendMessageResponse.status, sendMessageResponse.body);
  }
}
