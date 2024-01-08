# Starting situation
- Broker setup contains 3 brokers and three cruise control instances. 

# Happy cases
- The Kafka setup has three running brokers deployed. [1]
- The Kafka setup has one running operator pod. [2]
- The Kafka setup has one cruise control pod running. []

# Fault cases
- There are more than three Kafka brokers running. []
- There are less than three Kafka brokers running. []
- There is no operator pod running. []
- There are more than one operator pod running. []

# Test cases
| # | Test case | Desired outcome | Actual outcome |
| --- | --- | --- | --- |
| [case_1](case1_test.go) | Check the amount of running brokers in the kafka-brokers namespace. | There are three brokers running. |  |
| [case_2](case2_test.go) | Check the amount of running operator pods in the kafka-operator namespace. | There is one operator pod running. |  |
| [case_3](case3_test.go) | Check the amount of running cruise control pods in the kafka-brokers namespace. | There is one cruise control pod running. |  |

# Documentation Review
| # | Test case | Desired outcome |
| --- | --- | --- |
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of automatic queue deletion, including any configurable parameters and troubleshooting steps. |