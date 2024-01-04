package SSL_test

import (
    "io/ioutil"
    "testing"
    "time"
    "fmt"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestArtemisSSL(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "ArtemisSSL Suite")
}

var _ = Describe("Artemis SSL Connection", func() {
    Context("When inspecting the SSL certs directory", func() {
        It("should print the contents of the directory", func() {
            files, err := ioutil.ReadDir("/etc/ssl/certs")
            Expect(err).NotTo(HaveOccurred())

            for _, file := range files {
                fmt.Println(file.Name())
            }
        })
    })
    time.Sleep(1 * time.Minute)
})
