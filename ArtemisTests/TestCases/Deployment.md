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
|1| Check in cluster if 3 brokers are running in the correct namespace.| Broker setup (3 brokers) are present and running in the broker namespace.||
|2| Check in cluster if 1 operator pod is running in the correct namespace.| Operator pod is present and running in the operator namespace. ||
|3| Check if the statefulset is deployed in the cluster.| Statefulset is deployed and present.||
|4| Check if the securityfile is deployed in the cluster.| Securityfile is deployed and present. ||
|5| Check if the two certicates are present in the cluster. | Two certificates are present on the cluster, one in the broker namespace and one in the cert-manager namespace. ||
|6| Check if bundle (ca-bundle) is present on the cluster.| Bundle is present and synced true.||
|7| Check cluster if the root secret in the cert-manager namespace is present. | The root secret is present in the cert-manager namespace. ||
|8| Check cluster if the two secrets in the operator namespace are present. | Two secrets are present in the operator namespace.||
|9| Check cluster if the seven secrets in the broker namespace are present. | Seven secrets are present in the broker namespace.||
|10| Check in cluster if 1 trustmanager pod is running in the correct namespace.| Trustmanager pod is running in the cert-manager namespace.||
|11| Check if the broker setup has three bound PVCs (Persistent Volume Claims). | Three PVCs  are present and have the status bound.||
|12| Check if two clusterissuers are present.| Two clusterissuers are present and have the status ready is true. ||


# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | . | 