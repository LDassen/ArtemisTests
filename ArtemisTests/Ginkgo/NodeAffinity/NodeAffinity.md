# Starting situation
- Broker setup contains 3 brokers. 
- The amount of nodes (cluster dependent).
- That port 61619 is opened without SSL.

# Happy cases
- Brokers are configured to be deployed on different nodes to the fullest extent based on the load on a node. [1]

# Fault cases
- Node affinity rules are not followed or configured correctly, impacting broker pod placement on the nodes. [2]

# Test cases
| # | Test case | Desired outcome | Actual outcome |
| --- | --- | --- | --- |
| [case_1](case1_test.go) | Check if the brokers aren't running on one node. | The brokers get deployed on atleast 2 or more different nodes. | |
| [case_2](case2_test.go) | Turn off node affinity in configuration. | The brokers get deployed on different nodes (as much as possible). | |

# Documentation Review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of the node affinity (high availability), including any configurable parameters and troubleshooting steps. | 
