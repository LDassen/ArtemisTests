# Starting situation
- Broker setup contains 3 brokers. 
- There are no premade queues. 
- That port 61619 is opened without SSL.

# Happy cases
- When a message is sent to a queue that doesn't exist, the queue is created automatically. [case 1]
- Messages that are sent to queues remain in the queues until they are retrieved. [case 2]
- Messages sent to existing queues do not trigger autocreation of queue. [case 3]

# Fault cases
- A queue is not created when a message is sent to a queue that does not exist. [case 1] (Autocreation for queues must be false.)
- A message meant for a non-existent queue disappears. [case 4]

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
| [case_1](case1case2case3_test.go) | Publish a message to a non-existent queue. | The queue is created automatically and retains the new message. ||
| [case_2](case1case2case3_test.go) | Publish a message to a queue and retrieve this message after 10 minutes. | The message can be retrieved and is not deleted. ||
| [case_3](case1case2case3_test.go) | Publish a message to an already existing queue | The message is sent to the right queue and no other queue is created. ||
| 4 |  |  |  |
| 5 |  |  |  |

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of automatic queue creation, including any configurable parameters and troubleshooting steps. | 