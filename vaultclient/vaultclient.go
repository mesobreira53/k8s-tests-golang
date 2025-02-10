package vaultclient

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/vault/api"
)

// Attributes will be private
type VaultClient struct {
	address string
	token   string
}

type VaultStorage interface {
	VaultSyncSecret()
}

func NewVault(address string, token string) VaultClient {
	return VaultClient{address: address, token: token}
}

func (v VaultClient) VaultSyncSecret() {
	// Create a new Vault client configuration
	config := api.DefaultConfig()
	config.Address = v.address

	// Create the Vault client
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println("unable to create Vault client")
		log.Fatalf("unable to create Vault client: %v", err)
	}

	// Set the Vault token for authentication
	client.SetToken(v.token)

	// Define the secret data to store
	secretData := map[string]interface{}{
		"username": "vaultuser",
		"password": "vaultpassword",
	}

	// Write the secret to Vault at a specific path (e.g., secret/myapp)
	secretPath := "secret/data/myapp"
	_, err = client.KVv2("secret").Put(context.Background(), secretPath, secretData)
	if err != nil {
		log.Fatalf("unable to write secret: %v", err)
	}

	// Output a success message
	fmt.Printf("Successfully created secret at %s\n", secretPath)
}
