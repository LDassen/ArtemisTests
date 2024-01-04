# Starting situation
- Broker setup contains 3 brokers. 
- Predefined broker configuration is correct. 
- That port 61619 is opened without SSL.

# Happy cases
- Brokers can be deployed automatically when there is too much load for the current brokers. [1]
- Brokers can be removed automatically when there is little load for the current brokers. [2]
- Brokers can be deployed manually when there is too much load for the current brokers. [3]
- Brokers can be removed manually when there is too little load for the current brokers. [4]
- The broker number doesn't go lower than the minimum of three. [5]

# Fault cases
- Brokers cannot be deployed automatically when there is too much load for the current brokers. [1]
- Brokers cannot be removed automatically when there is too little load for the current brokers. [2]
- Brokers cannot be deployed manually when there is too much load for the current brokers. [3]
- Brokers cannot be removed manually when there is too little load for the current brokers. [4]
- The broker number does go lower than the minimum of three. [5]
- Failure to deploy automatically or manually due to resource restrictions. [6]
- Scaling up beyond max capacity. [6]

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
|1|Increase the load of the brokers so that a new broker is needed.|A new broker is automatically deployed.||
|2|(Do this with 4 or more brokers.)Decrease the load of the brokers so one broker is not needed anymore.|One broker is killed and the others pick up the load.||
|3|Go into the deployment file and add a broker. Then apply the file to the cluster.|The cluster adds a broker.||
|4|Go into the deployment file and remove a broker. Then apply this file to the cluster.|The cluster removes one broker.||
|5|Scale down to the minimum number of brokers (3). Try and reduce the load to see if downscaling happens or not.|The broker number does not go below three eventhough the load is significantly low.||
|6|Increase the load of the cluster that it reaches maximum capacity.|The system either rejects more deployments or gives warnings about the maximum capacity reached.||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the behavior of upscaling & downscaling, including any configurable parameters and troubleshooting steps. | 