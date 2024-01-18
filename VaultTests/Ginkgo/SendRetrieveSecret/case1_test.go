package SendRetrieveSecret_test

import (
    "github.com/hashicorp/vault/api"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("Vault Test with Kubernetes Job Configuration", func() {
    var client *api.Client
    var err error
    var token string

    BeforeEach(func() {
        // Initialize Vault client with TLS configuration
        config := &api.Config{
            Address: "https://dev-vault-1.dev-vault-1:8200",
        }
        config.ConfigureTLS(&api.TLSConfig{
            CACert:   "/etc/secret/ca.crt", // Path to the CA certificate
            Insecure: false,
        })

        client, err = api.NewClient(config)
        Expect(err).NotTo(HaveOccurred())

        // Hardcoded writer credentials (replace these with actual credentials)
        writerUsername := "writer"
        writerPassword := "Winter2018"

        // Authenticate with Vault using the writer credentials
        options := map[string]interface{}{
            "password": writerPassword,
        }
        path := "auth/userpass/login/" + writerUsername
        secret, err := client.Logical().Write(path, options)
        Expect(err).NotTo(HaveOccurred())
        Expect(secret).NotTo(BeNil())
        Expect(secret.Auth).NotTo(BeNil())

        token = secret.Auth.ClientToken
    })

    It("should store and retrieve a secret using writer credentials", func() {
        // Set the token for the client
        client.SetToken(token)

        secretPath := "application-vault/my-secret"
        secretData := map[string]interface{}{
            "key": "my-very-secure-value",
        }

        // Write a secret to the application-vault path
        _, err = client.Logical().Write(secretPath, secretData)
        Expect(err).NotTo(HaveOccurred())

        // Read the secret back to verify
        readSecret, err := client.Logical().Read(secretPath)
        Expect(err).NotTo(HaveOccurred())
        Expect(readSecret).NotTo(BeNil())
        Expect(readSecret.Data["data"].(map[string]interface{})["key"]).To(Equal("my-very-secure-value"))
    })
})
