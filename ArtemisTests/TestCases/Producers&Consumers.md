# Pre-requisites
- Multibroker of 3 pods is running.
- Login is working properly.
- Security is working properly.
- That port 61619 is opened without SSL.

# Happy cases
- The producers & consumers communicate with Activemq Artemis using the AMQP or CORE protocol. [1]
- A producer can send messages to a queue. [3]
- A consumer can read messages from a queue. [4]
- A queue only allows one producer to be connected. [5]
- A queue can have multiple consumers. [6]
- The producers & consumers get moved to a different broker if the broker they were connected to breaks or restarts. [9]

# Fault
- The producers & consumers cannot connect to Activemq Artemis to produce messages. [2]
- The producers cannot add messages to a queue. [2]
- The consumers cannot read messages from a queue. [2]
- Only one consumer is allowed per queue. [7]
- There are multiple producers connecting to the same queue. [8]
- Consumers & producers do not get reconnected to a different broker. [10]

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
|1|Connect with a producer & consumer Activemq Artemis using either the AMQP or CORE protocol.|The producer & consumer are connected.||
|2|Connect a producer & consumer on the same queue not using the AMQP or CORE protocols.|The producers & consumers cannot connect to Activemq Artemis.|| 
|3|Send a message to a queue using a producer.|The message can be seen in the queue.||
|4|Consume a message in queue that has messages.|It can be seen that the messages get consumed and acknowledged in the queue.||
|5|Connect with multiple producers to connect on one queue.|T.B.D.||
|6|Connect with multiple consumers to connect on one queue.|T.B.D.||
|7|T.B.D.|Only one consumer is allowed per queue.||
|8|T.B.D.|Multiple producers are connecting to one queue.||
|9|Broker with queues and its connected producers & consumers breaks/restarts. |The producers & consumers get reconnected (moved) to a different broker.||
|10|T.B.D.|The producers & consumers don't get reconnected (moved) to a different broker.||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | |