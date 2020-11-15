package providers

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault/keyvaultapi"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

func AzureVaultUrl(vaultName string) string {
	return fmt.Sprintf("https://%s.vault.azure.net", vaultName)
}

type AzureVaultProvider struct {
	vaultUrl string
	Client   keyvaultapi.BaseClientAPI
}

func NewAzureVaultProvider(vaultName string) (SecretsProvider, error) {
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, err
	}
	client := keyvault.New()
	client.Authorizer = authorizer
	return &AzureVaultProvider{vaultUrl: AzureVaultUrl(vaultName), Client: client}, nil
}

func (p *AzureVaultProvider) GetSecret(name string) (string, error) {
	resp, err := p.Client.GetSecret(context.Background(), p.vaultUrl, name, "")
	if err != nil {
		return "", err
	}
	return *resp.Value, err
}

func (p *AzureVaultProvider) GetKVSecrets(name string) (map[string]string, error) {
	return jsonKVSecrets(p, name)
}
