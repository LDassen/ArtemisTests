# Starting situation
- No brokers are running.
- Predefined broker configuration is correct. 
- Broker setup contains 3 brokers. 
- That port 61619 is opened without SSL.

# Happy cases
- Three brokers are successfully set up and are running. [1]

# Fault cases
- Multiple broker deployment fails. [2]
- Broker image cannot be pulled. [3]
- Running brokers are not equal to 3. [4]

# Test cases
| # | Test case | Desired outcome | Actual outcome |
| --- | --- | --- | --- |
| [case_1](case1_test.go) | Deploy broker configuration. | Three brokers are deployed successfully and are running. | |
| [case_2](case2_test.go) | Deploy the broker configuration file in a namespace that doesn't exist. | Brokers cannot be deployed and gets a "namespace not found" error. | |
| [case_3](case3_test.go) | T.B.D. | Brokers cannot be deployed due to image pull error. | |
| [case_4](case4_test.go) | T.B.D. |  | |

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of the multiple broker setup, including any configurable parameters and troubleshooting steps. | 