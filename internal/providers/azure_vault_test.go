package providers

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault/keyvaultapi"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeAzureClient struct {
	keyvaultapi.BaseClientAPI
	data map[string]string
}

func newFakeAzureClient(data map[string]string) *fakeAzureClient {
	return &fakeAzureClient{data: data}
}

func (f *fakeAzureClient) GetSecret(ctx context.Context, vaultBaseURL string, secretName string, secretVersion string) (keyvault.SecretBundle, error) {
	if val, found := f.data[secretName]; found {
		return keyvault.SecretBundle{Value: &val}, nil
	}
	return keyvault.SecretBundle{}, EmptySecretError
}

func TestAzureVaultUrl(t *testing.T) {
	res := AzureVaultUrl("test")
	assert.Equal(t, "https://test.vault.azure.net", res)
}

func TestAzureVaultProvider_GetSecret(t *testing.T) {
	p := AzureVaultProvider{Client: newFakeAzureClient(TestMap)}
	res, err := p.GetSecret(TestKey)
	assert.NoError(t, err)
	assert.Equal(t, TestValue, res)
}

func TestAzureVaultProvider_GetSecretErr(t *testing.T) {
	p := AzureVaultProvider{Client: newFakeAzureClient(map[string]string{})}
	_, err := p.GetSecret(TestKey)
	assert.Error(t, err)
}

func TestAzureVaultProvider_GetKVSecrets(t *testing.T) {
	p := AzureVaultProvider{Client: newFakeAzureClient(TestMapJson)}
	res, err := p.GetKVSecrets(TestKey)
	assert.NoError(t, err)
	assert.Equal(t, TestSecretValue, res[TestSecretKey])
}

func TestAzureVaultProvider_GetKVSecretsErr(t *testing.T) {
	client := newFakeAzureClient(map[string]string{})
	p := AzureVaultProvider{
		Client: client,
	}
	_, err := p.GetKVSecrets(TestKey)
	assert.Error(t, err)

	// non-json value
	client.data = TestMap
	_, err = p.GetKVSecrets(TestKey)
	assert.Error(t, err)
}

func TestNewAzureVaultProvider(t *testing.T) {
	_, err := NewAzureVaultProvider("")
	assert.NoError(t, err)
}