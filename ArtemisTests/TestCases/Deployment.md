# Starting situation
- Broker setup contains 3 brokers. 
- For every broker there is a linking persistent volume.
- That port 61619 is opened without SSL.

# Happy cases
- x

# Fault
- x

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
|1| Check in cluster if 3 brokers are running in the correct namespace.|||
|2| Check in cluster if 1 operator pod is running in the correct namespace.| ||
|3| Check if the statefulset is deployed in the cluster.| ||
|4| Check if the securityfile is deployed in the cluster.| ||
|5| Check if the two certicates are present in the cluster. | ||
|6| Check if bundle (ca-bundle) is present on the cluster.| ||
|7| . | ||
|8| ..| ||

# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | . | 