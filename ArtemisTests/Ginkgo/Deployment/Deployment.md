# Starting situation
- Broker setup contains 3 brokers. 
- For every broker there is a linking persistent volume.
- That port 61619 is opened without SSL.

# Test cases
|#|Test case|Desired outcome|Actual outcome|
|---|---|---|---|
|[case1](case1_test.go)| Check in cluster if 3 brokers are running in the correct namespace.| Broker setup (3 brokers) are present and running in the broker namespace.||
|[case2](case2_test.go)| Check in cluster if 1 operator pod is running in the correct namespace.| Operator pod is present and running in the operator namespace. ||
|[case3](case3_test.go)| Check if the statefulset is deployed in the cluster.| Statefulset is deployed and present.||
|[case4](case4_test.go)| Check if the securityfile is deployed in the cluster.| Securityfile is deployed and present. ||
|[case5](case5_test.go)| Check if the two certicates are present in the cluster. | Two certificates are present on the cluster, one in the broker namespace and one in the cert-manager namespace. ||
|[case6](case6_test.go)| Check if bundle (ca-bundle) is present on the cluster.| Bundle is present and synced true.||
|[case7](case7_test.go)| Check cluster if the root secret in the cert-manager namespace is present. | The root secret is present in the cert-manager namespace. ||
|[case8](case8_test.go)| Check cluster if the two secrets in the operator namespace are present. | Two secrets are present in the operator namespace.||
|[case9](case9_test.go)| Check cluster if the seven secrets in the broker namespace are present. | Seven secrets are present in the broker namespace.||
|[case10](case10_test.go)| Check in cluster if 1 trustmanager pod is running in the correct namespace.| Trustmanager pod is running in the cert-manager namespace.||
|[case11](case11_test.go)| Check if the broker setup has three bound PVCs (Persistent Volume Claims). | Three PVCs  are present and have the status bound.||
|[case12](case12_test.go)| Check if two clusterissuers are present.| Two clusterissuers are present and have the status ready is true. ||


# Documentation review
| # | Test case | Desired outcome |
| --- | --- | --- | 
| # | Review documentation in ADO WIKI. | Confirm that the documentation accurately reflects the starting situation of the complete deployment. | 