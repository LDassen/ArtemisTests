import { check } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  const url = 'http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619';
  const queueName = 'exampleQueueCore';
  const message = 'Hello, Core!';

  const response = check(
    ws.connect(url, {}, (socket) => {
      // Send message
      socket.send(JSON.stringify({ queueName, message }));

      // Receive message
      const receivedMessage = socket.recv();
      check(receivedMessage, {
        'Core Message Received Successfully': (r) => r !== null,
      });
    }),
    { 'WebSocket Connection Established': (r) => r && r.code === 101 }
  );
}
