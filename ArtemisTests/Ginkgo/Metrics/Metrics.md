# Starting situation
- Artemis setup contains 3 brokers. 
- That port 61619 is opened without SSL.

# Happy cases
- The Artemis setup produces metrics. [1]

# Fault cases
- No metrics are produced. [1]

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
| [case_1](case1_test.go) | Go into the scrape port and see if there are metrics being produced. | The Artemis setup produces metrics. ||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of automatic queue creation, including any configurable parameters and troubleshooting steps. | 