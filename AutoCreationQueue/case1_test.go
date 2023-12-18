package AutoCreationQueue_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands"
	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common"
	"os/exec"
)

var _ = Describe("ActiveMQ Artemis Brokers Pods", func() {
	const namespace = "activemq-artemis-brokers"

	Context("When querying ActiveMQ Artemis Brokers pods", func() {
		It("Should retrieve a list of pods", func() {
			cmd := exec.Command("kubectl", "get", "pods", "-n", namespace)
			session, err := commands.Execute(cmd)
			Expect(err).NotTo(HaveOccurred())

			Expect(session.ExitCode()).To(Equal(0), "kubectl command failed. Output:\n%s", session.Out.Contents())

			// Add more assertions as needed based on the output
			// For example, you can use Gomega matchers to check the expected output or count of pods.
		})
	})
})
