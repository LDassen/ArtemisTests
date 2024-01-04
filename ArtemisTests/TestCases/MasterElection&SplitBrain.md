# Starting situation
- Broker setup contains 3 brokers. 
- That port 61619 is opened without SSL.

# Happy cases
- Broker configuration effectively prevents split brain scenarios. [ ]
- The system successfully resolves split brain situations and restores consistency. [ ]
- Master election occurs consistenly. [ ]

# Fault cases
- Split brain occurs, leading to divergent states in different parts of the cluster. [ ]
- The system fails to recover and heal after resolving network partitions. [ ]
- Using a single primary/backup pair, so there is no mitigation against split brain. [ ]

# Test cases
| # | Test case | Desired outcome | Actual outcome |
| --- | --- | --- | --- |
| 1 | Configure and test quorum voting within the broker setup. | The brokers get deployed and the replicas use the voting quorum to decide (how long it takes) to takeover the live server or just stay a backup server. |  |
| 2 | Using a single primary/backup pair | There is no mitigation against split brain, it's not possible for the brokers to vote so both pods will become a live server.  |  |
| 3 |  Intentionally create network partition. | The system successfully detects and resolves the split brain scenario, restoring consistency. | |
| 4 | Simulate split brain with (bad) configurations. | The brokers get isolated from the network. |  |
| 5 | Test quorum voting with different timeout configurations. | The brokers still behaves according to the configured quorum voting timeout, preventing and resolving split brain scenarios. | |
| 6 | Continuously kill the master broker. | A new pod will be quorum voted to become a master and if you stop killing the previous master broker, it will become a slave pod. | |

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of network isolation, master election and quorum voting, including any configurable parameters and troubleshooting steps. | 