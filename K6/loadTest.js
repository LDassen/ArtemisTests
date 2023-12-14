import { check } from 'k6';
import { AMQP, Core } from 'k6/x/amqp';
import { SharedArray } from 'k6/data';

// Define the Artemis broker connection details
const AMQP_CONFIG = {
  url: 'amqp://cgi:cgi@ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619',
};

const CORE_CONFIG = {
  url: 'tcp://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619',
};

// Function to send and receive messages using AMQP
function sendAndReceiveAMQPMessage(session) {
  const queueName = 'exampleQueueAMQP';
  const message = 'Hello, AMQP!';

  // Send message
  const sendResult = session.send(queueName, message);
  check(sendResult, {
    'AMQP Message Sent Successfully': (r) => r === 0,
  });

  // Receive message
  const receiveResult = session.receive(queueName);
  check(receiveResult, {
    'AMQP Message Received Successfully': (r) => r !== null,
  });
}

// Function to send and receive messages using Core
function sendAndReceiveCoreMessage(session) {
  const queueName = 'exampleQueueCore';
  const message = 'Hello, Core!';

  // Send message
  const sendResult = session.send(queueName, message);
  check(sendResult, {
    'Core Message Sent Successfully': (r) => r === 0,
  });

  // Receive message
  const receiveResult = session.receive(queueName);
  check(receiveResult, {
    'Core Message Received Successfully': (r) => r !== null,
  });
}

export const options = {
  vus: 10,
  duration: '30s',
};

const amqpSessions = new SharedArray('amqpSessions', function () {
  const sessions = [];
  for (let i = 0; i < options.vus; i++) {
    const amqpSession = new AMQP(AMQP_CONFIG);
    sessions.push(amqpSession);
  }
  return sessions;
});

const coreSessions = new SharedArray('coreSessions', function () {
  const sessions = [];
  for (let i = 0; i < options.vus; i++) {
    const coreSession = new Core(CORE_CONFIG);
    sessions.push(coreSession);
  }
  return sessions;
});

export default function () {
  for (let i = 0; i < options.vus; i++) {
    sendAndReceiveAMQPMessage(amqpSessions[i]);
    sendAndReceiveCoreMessage(coreSessions[i]);
  }
}
