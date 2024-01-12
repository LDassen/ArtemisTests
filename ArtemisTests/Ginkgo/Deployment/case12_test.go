package Deployment_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Check ClusterIssuer Existence", func() {
	It("should ensure ClusterIssuers exist in the YAML files associated with certificates", func() {
		// Specify the path to your YAML file containing certificates and ClusterIssuers
		yamlFilePath := filepath.Join(homedir.HomeDir(), "path/to/your/file.yaml")

		// Read the YAML file
		yamlContent, err := ioutil.ReadFile(yamlFilePath)
		Expect(err).To(BeNil(), "Error reading YAML file: %v", err)

		// Unmarshal YAML content into a map
		var yamlData map[string]interface{}
		err = yaml.Unmarshal(yamlContent, &yamlData)
		Expect(err).To(BeNil(), "Error unmarshalling YAML content: %v", err)

		// Check if the YAML content contains certificates and ClusterIssuers
		certificates, foundCertificates := yamlData["Certificates"].(map[string]interface{})
		Expect(foundCertificates).To(BeTrue(), "No certificates found in the YAML file.")

		for certName, certData := range certificates {
			// Assuming the YAML structure, adjust it based on your actual structure
			clusterIssuerName, found := certData.(map[string]interface{})["ClusterIssuer"].(string)
			Expect(found).To(BeTrue(), "ClusterIssuer not found for certificate %s", certName)package Deployment_test

			import (
				"context"
				"fmt"
				. "github.com/onsi/ginkgo/v2"
				. "github.com/onsi/gomega"
				"k8s.io/client-go/kubernetes"
				"k8s.io/client-go/rest"
				metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
				certmanagerv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
			)
			
			var _ = Describe("Check ClusterIssuers Existence", func() {
				It("should ensure ClusterIssuers exist in the specified namespace", func() {
					config, err := rest.InClusterConfig()
					Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)
			
					clientset, err := kubernetes.NewForConfig(config)
					Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)
			
					namespace := "cert-manager"
			
					// List all ClusterIssuers in the namespace
					clusterIssuersList, err := clientset.CertmanagerV1().ClusterIssuers().List(context.TODO(), metav1.ListOptions{})
					Expect(err).To(BeNil(), "Error listing ClusterIssuers: %v", err)
			
					// Names of ClusterIssuers to find
					clusterIssuerNames := []string{"amq-ca-issuer", "amq-selfsigned-cluster-issuer"}
			
					// Check each ClusterIssuer's existence and readiness
					for _, clusterIssuerName := range clusterIssuerNames {
						found := false
						for _, ci := range clusterIssuersList.Items {
							if ci.Name == clusterIssuerName {
								found = true
								fmt.Printf("ClusterIssuer '%s' found in namespace '%s'\n", ci.Name, namespace)
			
								// Perform additional checks if needed
			
								// Check the conditions
								Expect(ci.Status.Conditions).To(HaveLen(1), "Expected ClusterIssuer to have one condition.")
								Expect(ci.Status.Conditions[0].Type).To(Equal(certmanagerv1.ConditionReady), "Expected ClusterIssuer condition to be Ready.")
								Expect(ci.Status.Conditions[0].Status).To(Equal(certmanagerv1.ConditionTrue), "Expected ClusterIssuer condition status to be True.")
								break
							}
						}
						Expect(found).To(BeTrue(), "ClusterIssuer '%s' not found in namespace '%s'", clusterIssuerName, namespace)
					}
				})
			})
			

			// Check if the ClusterIssuer exists in the YAML content
			_, found = yamlData["ClusterIssuers"].(map[string]interface{})[clusterIssuerName]
			if !found {
				// ClusterIssuer not found
				fmt.Printf("ClusterIssuer '%s' not found in the YAML file\n", clusterIssuerName)
				Expect(found).To(BeTrue(), "Expected ClusterIssuer to be present, but it is missing.")
			} else {
				// ClusterIssuer found
				fmt.Printf("ClusterIssuer '%s' found in the YAML file\n", clusterIssuerName)
				Expect(found).To(BeTrue(), "Expected ClusterIssuer to be present, but it is missing.")
			}
		}
	})
})
