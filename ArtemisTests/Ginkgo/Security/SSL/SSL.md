# Starting situation
- Artemis setup contains 3 brokers. 
- That port 61617 is opened with SSL.

# Happy cases
- The Artemis setup uses SSL for its communication. [1]
- SSL is setup correctly via the bundle, i.e. every namespace has a bundle present. [2]

# Fault cases
- Communication in the Artemis setup uses SSL but can be accessed without SSL. [3]
- It is not possible to connect to the Artemis setup with a producer/consumer which has the correct configmap (bundle) mounted. []

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
| [case_1](case1_test.go) | Setup a producer that sends a message to a queue. The producer should have the bundle mounted. | The message is produced on the queue. ||
| [case_2](case2_test.go) | Check every namespace in the cluster to see if the bundle configmap is present. | Every namespace in the cluster has the bundle. ||
| [case_3](case3_test.go) | Setup a producer that sends a message to a queue. The producer should not have the bundle mounted. | The message is not produced inside the queue. ||
| [case_4](case4_test.go) | Setup a producer that sends a message to a queue. The producer should have the bundle mounted. | The message is not produced inside the queue. ||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of automatic queue creation, including any configurable parameters and troubleshooting steps. | 