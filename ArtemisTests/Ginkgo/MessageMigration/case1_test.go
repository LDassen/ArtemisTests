package MessageMigration_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	//apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = ginkgo.Describe("Kubernetes Apply CRD Test", func() {
	var clientset *kubernetes.Clientset

	ginkgo.BeforeEach(func() {
		// Set up the Kubernetes client
		config, err := rest.InClusterConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		clientset, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("Should reapply a CRD file for CustomResourceDefinition", func() {
		fileName := "ex-aaoMM.yaml"

		// Read the file
		filePath, err := filepath.Abs(filepath.Join(homedir.HomeDir(), "path/to", fileName))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		fileBytes, err := ioutil.ReadFile(filePath)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Decode the YAML manifest
		decode := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(fileBytes), 1024)
		var crd apiextv1.CustomResourceDefinition
		err = decode.Decode(&crd)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Try to apply the CRD
		err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			_, updateErr := clientset.ApiextensionsV1().CustomResourceDefinitions().Update(context.TODO(), &crd, metav1.UpdateOptions{})
			if updateErr == nil {
				fmt.Println("CRD re-applied successfully!")
				return nil
			}
			fmt.Println("[ERROR] Error reapplying CRD:", updateErr)
			return updateErr
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to reapply CRD")
	})
})
