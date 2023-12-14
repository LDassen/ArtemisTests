import { check, sleep } from 'k6';
import http from 'k6/http';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  const url = 'http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619';
  const queueName = 'TEST';
  const message = 'Hello, Core!';
  const username = 'cgi';
  const password = 'cgi';

  const payload = JSON.stringify({
    queueName,
    message,
  });

  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Basic ${customEncodeBase64(`${username}:${password}`)}`,
  };

  const response = http.post(`${url}/send-receive-endpoint`, payload, { headers });
  console.log(response.status, response.body);
  console.log(`HTTP Request: ${JSON.stringify({ url: `${url}/send-receive-endpoint`, payload, headers }, null, 2)}`);
  console.log(`HTTP Response: ${JSON.stringify(response, null, 2)}`);
  console.log(`HTTP Response Status Code: ${response.status}`);

  check(response, {
    'HTTP Request Successful': (r) => r.status === 200,
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
