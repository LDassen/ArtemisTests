# Starting situation
- Broker setup contains 3 brokers. 
- For every broker there is a linking persistent volume.
- That port 61619 is opened without SSL.

# Happy cases
- Messages and queue information will remain on the PVs when the broker breaks or restarts. [1]
- The PV will retain its messages under high load. [2, 3]
- A replacement persistent volume will be filled with the information from the previous one. [4]
- Messages with an expiration time will disappear from the queue and PV when they are expired. [5]
- Persistent volumes have data distribution and sharing. []

# Fault cases
- Messages and queue information are not retained in the PVs when a broker breaks or restarts. [1]
- The PV does not retain all messages under high load. [2, 3]
- A replacement persistent volume is not filled with its predecessors information. [4]
- Messages with an expiration date do not disappear from the queue and PV when they are expired. [5]

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
|1|Send a message to a queue. Restart or kill the broker to see whether the PV will retain its data.|The PV retains the data.||
|2|Send in a large number of messages.|The PV retains all messages.||
|3|Simultaneously send and consume messages from the same queue.|The PV stays consistent and ensures data integrity.||
|4|Kill a broker.|The persistent volume should not be removed and the data inside the persistent volume remains intact.||
|5|Send messages with expiration dates to a queue and wait for them to expire.|The queue and PV should not have these messages anymore.||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of message queue retetion, including any configurable parameters and troubleshooting steps. | 