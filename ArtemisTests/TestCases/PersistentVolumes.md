# Starting situation
- Broker setup contains 3 brokers. 
- That port 61619 is opened without SSL.

# Happy cases
- The persistent volumes store the queue information and messages. [1]
- A persistent volume will remain when a broker breaks or restarts. [2]
- The persistent volume claim will not be removed when a broker breaks or restarts. [3]
- PV can be resized while the broker is running without any issues. [4]
- PVs will expand their storage when getting close to being full. [5]
- Replacing a PV when under heavy load is possible. [6]

# Fault cases
- The persistent volumes do not properly store the queue information and messages. [1]
- Persistent volumes disappear when a broker breaks or restarts. [2]
- Persistent volume claims disappear after a broker breaks or restarts. [3]
- Resizing a PV while the broker is running will result in issues. [4]
- PVs do not expand when getting close to being full. [5]
- Replacing a PV when under heavy load is not possible. [6]

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
|1|Send a message to a queue. Check if the PV has this message.|The PV has the message stored.||
|2|Kill a broker.|The persistent volume does not disappear and is used for the new broker.||
|3|Remove a persistent volume.|A new persistent volume is created and it is filled with the information of the old one. *This also means that the persistent volume claim owner is not the broker.*||
|4|Resize the PV while the broker is running.|PV handles resizing without data loss and adapts to changes.||
|5|Simulate a PV that is almost at full capacity.|The PV will expand itself.||
|6|Simulate hevay load for the brokers and remove a PV.|The replacement happens without any problems.||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of persistent volumes, including any configurable parameters and troubleshooting steps. | 