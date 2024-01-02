# Starting situation
- Broker setup contains 3 brokers. 
- For every broker there is a linking persistent volume.
- That port 61619 is opened without SSL.

# Happy cases
- Messages get migrated successfully from one broker to another broker. [1] 
- An drainer pod will automatically start and handle the message migration. [1] 
- The drainer pod will automatically dissapear after message migration. [1] 

# Fault cases
- During the time the drainer pod migrates messages, the running pods are 3 (brokers) + 1 (drainer). [1]
- Message migration failed from one broker to another broker. [2]
- The drainer pod will not distribute load to multiple brokers when migrating messages. [2]
- After message migration failure the persistent volume of the broker will not be deleted. [2]

# Test cases
| # | Test case | Desired outcome | Actual outcome |
| --- | --- | --- | --- |
| [case_1](case1_test.go)  | Delete a broker pod with queues and messages on it. | A drainer pod will be automatically created and migrates all messages from that broker to another running broker. | |
| [case_2](case2_test.go)  | Don't enable message migration within the CRD. | No drainer pod will be running after a broker gets killed and message migration fails. | |

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of automatic message migration, including any configurable parameters and troubleshooting steps. | 