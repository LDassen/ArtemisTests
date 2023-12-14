import { check, sleep } from 'k6';
import http from 'k6/http';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  const url = 'http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619';
  const queueName = 'TESTKUBE';
  const message = 'Hello, Core!';
  const user = 'amq';
  const password = 'amq';

  // The payload structure might depend on your specific use case and Artemis configuration
  const payload = JSON.stringify({
    body: message,
  });

  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Basic ${customEncodeBase64(`${user}:${password}`)}`,
  };

  // Perform a HEAD or GET request on the queue to get necessary information
  const responseQueue = http.get(`${url}/queues/${queueName}`, { headers });

  // Extract necessary URLs from the response headers
  const msgCreateUrl = responseQueue.headers['msg-create'];
  const msgCreateWithIdUrl = responseQueue.headers['msg-create-with-id'];

  // Perform a POST request to create a message in the queue
  const responseCreateMsg = http.post(msgCreateUrl, payload, { headers });

  console.log(responseCreateMsg.status, responseCreateMsg.body);
  console.log(`HTTP Request: ${JSON.stringify({ url: msgCreateUrl, payload, headers }, null, 2)}`);
  console.log(`HTTP Response: ${JSON.stringify(responseCreateMsg, null, 2)}`);
  console.log(`HTTP Response Status Code: ${responseCreateMsg.status}`);

  // Check the response from the POST request
  check(responseCreateMsg, {
    'Message Creation Successful': (r) => r.status === 200,
    // You may need to adjust the check based on the actual response structure
    'Core Message Received Successfully': (r) => r.json('receivedMessage') !== null,
  });

  // Sleep for a short duration between requests (adjust as needed)
  sleep(1);
}

function customEncodeBase64(str) {
  const keyStr = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=';
  let output = '';
  let chr1, chr2, chr3, enc1, enc2, enc3, enc4;
  let i = 0;

  while (i < str.length) {
    chr1 = str.charCodeAt(i++);
    chr2 = str.charCodeAt(i++);
    chr3 = str.charCodeAt(i++);

    enc1 = chr1 >> 2;
    enc2 = ((chr1 & 3) << 4) | (chr2 >> 4);
    enc3 = ((chr2 & 15) << 2) | (chr3 >> 6);
    enc4 = chr3 & 63;

    if (isNaN(chr2)) {
      enc3 = enc4 = 64;
    } else if (isNaN(chr3)) {
      enc4 = 64;
    }

    output = output +
      keyStr.charAt(enc1) + keyStr.charAt(enc2) +
      keyStr.charAt(enc3) + keyStr.charAt(enc4);
  }

  return output;
}
