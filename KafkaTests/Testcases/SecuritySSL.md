# Starting situation
- Kafka setup contains 3 brokers. 
- That port 9094 is open with SSL.

# Happy cases
- The Kafka setup uses SSL for its communication. [1]
- SSL is set up correctly. [2]

# Fault cases
- Communication in the KAfka setup uses SSL but can be accessed without SSL. [3]
- It is not possible to connect to the Kafka setup with a producer/consumer which has the correct secret mounted. []

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
| 1 | Setup a producer that sends a message to a queue. The producer should have the secret mounted. | The message is produced on the queue. ||
| 2 | Check both certificates and secrets in the cluster to see if they are present. | They are present. ||
| 3 | Setup a producer that sends a message to a queue. The producer should not have the secret mounted. | The message is not produced inside the queue. ||
| 4 | Setup a producer that sends a message to a queue. The producer should have the secret mounted. | The message is not produced inside the queue. ||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of automatic queue creation, including any configurable parameters and troubleshooting steps. | 