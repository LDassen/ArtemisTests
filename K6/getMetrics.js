import { check, sleep } from 'k6';
import http from 'k6/http';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  const url = 'http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:8161';  // Change the port to 8161
  const metricsEndpoint = '/metrics';  // Endpoint to retrieve metrics

  const response = http.get(`${url}${metricsEndpoint}`);

  check(response, {
    'HTTP Request Successful': (r) => r.status === 200,
    'Metrics Retrieved Successfully': (r) => r.body.includes('your_metric_keyword'),  // Replace 'your_metric_keyword' with an actual keyword from your metrics
  });

  // Log the metrics to the console
  console.log(`Metrics: ${response.body}`);

  // Sleep for a short duration between requests (adjust as needed)
  sleep(1);
}
