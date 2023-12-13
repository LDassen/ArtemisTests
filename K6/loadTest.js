import { exec } from 'k6/execution';
import { check } from 'k6';

export let options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  // Replace the URL, username, and password with your Artemis broker endpoint and credentials
  let artemisURL = 'http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619';
  let username = 'cgi';
  let password = 'cgi';
  let messageCount = 100;

  // Execute the Artemis CLI tool command
  let command = `./artemis producer --user ${username} --password ${password} --url ${artemisURL} --message-count ${messageCount}`;
  let result = exec(command);

  // Check if the command was successful (return code 0)
  check(result, {
    'command executed successfully': (r) => r.code === 0,
  });
}
