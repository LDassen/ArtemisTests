# Starting situation
- Broker setup contains 3 brokers and three cruise control instances. 

# Happy cases
- The Vault setup has five running instances deployed. [case 1]
- The Vault setup has one running operator pod. [case 2]
- The Vault setup has one vault configurer pod running. [case 3]
- The Vault setup has the statefulset deployed. [case 4]
- The Vault setup has five secrets deployed. [case 5]
- The Vault setup has a role defined. [case 6]
- The Vault setup has a rolebinding defined. [case 7]
- The Vault setup has a clusterrolebinding defined. [case 8]

# Test cases
| # | Test case | Desired outcome | Actual outcome |
| --- | --- | --- | --- |
| [case_1](case1_test.go) | Check the amount of running brokers in the vault namespace. | There are five instances running. |  |
| [case_2](case2_test.go) | Check the amount of running operator pods in the operator namespace. | There is one operator pod running. |  |
| [case_3](case3_test.go) | Check the amount of running vault configurer pods in the vault namespace. | There is one vault configurer pod running. |  |
| [case_4](case4_test.go) | Check if the statefulset is deployed in the cluster.| The statefulset is deployed and present.||
| [case_5](case5_test.go) | Check if the five secrets are present in the vault namespace. | There are five secrets present. |  |
| [case_6](case6_test.go) | Check if the role is defined. | The role is present and defined. ||
| [case_7](case7_test.go) | Check if the rolebinding is defined. | The rolebinding is present and defined.||
| [case_8](case8_test.go) | Check if the clusterrolebinding is defined. | The clusterrolebinding is present and defined. ||

# Documentation Review
Confirm that the documentation accurately reflects the behavior of automatic queue deletion, including any configurable parameters and troubleshooting steps.

# To do
case 6
case 7
case 8