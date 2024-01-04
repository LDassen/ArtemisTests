# Starting situation
- Broker setup contains 3 brokers. 
- That port 61619 is opened without SSL.

# Happy cases
- The Artemis setup has logging enabled. [1]
- Logging is received by ElasticSearch. [2]
- Loglevels can be changed when needed. [3]
- When a problem occurs, the log output will show the problem. [4]

# Fault cases
- Logging is not enabled for the Artemis setup. [1]
- ElasticSearch does not receive logging from Artemis. [2]
- Loglevels cannot be changed. [3]
- When a problem occurs, the log output does not show the problem. [4]

# Unknowns
- Adjusting log level might not be possible in Azure.

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
|1|Check the Artemis broker output.|The output shows logs.||
|2|Check the logging dashboard to see if ElasticSearch received any logs from Artemis.|ElasticSearch has received logs.||
|3|Change the loglevel from (for example) INFO to DEBUG.|The log output is in DEBUG instead of INFO.||
|4|Create a problem (for example, kill a broker).|The logs reflect their output to the problem.||

# Documentation Review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of logging, including any configurable parameters and troubleshooting steps. | 