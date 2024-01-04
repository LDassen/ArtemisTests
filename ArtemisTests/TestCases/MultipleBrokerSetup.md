# Starting situation
- No brokers are running.
- Predefined broker configuration is correct. 
- Broker setup contains 3 brokers. 

# Happy cases
- Three brokers are successfully set up and are running. [1]

# Fault cases
- Multiple broker deployment fails. [2]
- Broker image cannot be pulled. [3]
- Running brokers are not equal to 3. [4]

# Test cases
| # | Test case | Desired outcome | Actual outcome |
| --- | --- | --- | --- |
| 1 | Deploy broker configuration. | The ex-aao file is successfully deployed. | |
| 2 | Deploy the broker configuration file in a namespace that doesn't exist. | Brokers cannot be deployed and gets a "namespace not found" error. | |
| 3 | T.B.D. | Brokers cannot be deployed due to image pull error. | |
| 4 | T.B.D. |  | |

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of the multiple broker setup, including any configurable parameters and troubleshooting steps. | 