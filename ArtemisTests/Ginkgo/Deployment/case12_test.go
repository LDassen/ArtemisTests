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
			Expect(found).To(BeTrue(), "ClusterIssuer not found for certificate %s", certName)

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
