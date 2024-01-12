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
| [case_1](case1_test.go) | Setup a producer that sends a message to a queue. The producer should have the secret mounted. | The message is produced on the queue. ||
| [case_2](case2_test.go) | Check both certificates and secrets in the cluster to see if they are present. | They are present. ||
| [case_3](case3_test.go) | Setup a producer that sends a message to a queue. The producer should not have the secret mounted. | The message is not produced inside the queue. ||
| Needs to be done still. [case_4](case4_test.go) | Setup a producer that sends a message to a queue. The producer should have the secret mounted. | The message is not produced inside the queue. ||

*Case 3 and 4 will be put in another folder called SSL because they do not need to have the configmap mounted.*

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of automatic queue creation, including any configurable parameters and troubleshooting steps. | 