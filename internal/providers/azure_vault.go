package providers

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault/keyvaultapi"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"os"
	"strings"
)

func AzureVaultUrl(vaultName string) string {
	return fmt.Sprintf("https://%s.vault.azure.net", vaultName)
}

type AzureVaultProvider struct {
	vaultUrl string
	Client   keyvaultapi.BaseClientAPI
}

func NewAzureVaultProvider(vaultName string) (SecretsProvider, error) {
	authorizer, err := getAuthorizer()
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

func getAuthorizer() (autorest.Authorizer, error) {
	settings, err := auth.GetSettingsFromEnvironment()
	if err != nil {
		return nil, err
	}
	settings.Values[auth.Resource] = strings.TrimSuffix(settings.Environment.KeyVaultEndpoint, "/")

	// based on Azure SDK EnvironmentSettings.GetAuthorizer()
	//1.Client Credentials
	if c, e := settings.GetClientCredentials(); e == nil {
		return c.Authorizer()
	}

	//2. Client Certificate
	if c, e := settings.GetClientCertificate(); e == nil {
		return c.Authorizer()
	}

	//3. Username Password
	if c, e := settings.GetUsernamePassword(); e == nil {
		return c.Authorizer()
	}

	// 4. MSI
	if _, present := os.LookupEnv("AZURE_USE_MSI"); present {
		return settings.GetMSI().Authorizer()
	}

	// 5. CLI
	return auth.NewAuthorizerFromCLIWithResource(settings.Values[auth.Resource])
}
