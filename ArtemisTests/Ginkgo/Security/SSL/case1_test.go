package SSL_test

import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "net"
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestArtemisSSL(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "ArtemisSSL Suite")
}

var _ = Describe("Artemis SSL Connection", func() {
    Context("When connecting to Artemis with SSL", func() {
        It("should successfully connect", func() {
            caCert, err := ioutil.ReadFile("/etc/ssl/certs/ca-bundle.crt")
            Expect(err).NotTo(HaveOccurred())

            caCertPool := x509.NewCertPool()
            caCertPool.AppendCertsFromPEM(caCert)

            config := &tls.Config{
                RootCAs: caCertPool,
            }

            conn, err := tls.Dial("tcp", "<ARTEMIS_HOST>:<ARTEMIS_PORT>", config)
            Expect(err).NotTo(HaveOccurred())
            defer conn.Close()

            // Perform further checks if needed
        })
    })
})
