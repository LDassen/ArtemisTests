# Starting situation
- Broker setup contains 3 brokers. 
- There are no premade queues. 
- That port 61619 is opened without SSL.

# Happy cases
- Queues are removed in max 30 seconds when they have no messages or consumers connected. [1]

# Fault cases
- Queue deletion does not happen when there are no messages or consumers connected. [1] (Autodeletion for queues must be false.)
- Queue deletion occurs when messages and or consumers are connected. []
- Queue deletion occurs when messages are still present but no producer or consumer are connected. []

# Test cases
| # | Test case | Desired outcome | Actual outcome |
| --- | --- | --- | --- |
|[case_1](case1_test.go) | Send a message to a queue and retrieve that message. | The queue should delete after max 2 minutes (considering delays in testing). |  |
| 2 |  |  |  |
| 3 |  |  |  |

# Documentation Review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of automatic queue deletion, including any configurable parameters and troubleshooting steps. | 