package providers

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestSecretPath = "test"
)

type fakeVaultClient struct {
	secrets map[string]*api.Secret
}

func NewFakeVaultClient() *fakeVaultClient {
	return &fakeVaultClient{
		secrets: map[string]*api.Secret{},
	}
}

func (f *fakeVaultClient) SetSecret(path string, s *api.Secret) {
	f.secrets[path] = s
}

func (f *fakeVaultClient) Read(path string) (*api.Secret, error) {
	if val, found := f.secrets[path]; found {
		return val, nil
	}
	return nil, EmptySecretError
}

func TestNewHashicorpVaultProvider(t *testing.T) {
	p, err := NewHashicorpVaultProvider()
	assert.NoError(t, err)
	assert.NotNil(t, p)

}

func TestHashicorpVaultProvider_GetSecret(t *testing.T) {

	fakeClient := NewFakeVaultClient()

	fakeClient.SetSecret(TestSecretPath, &api.Secret{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				TestSecretKey: TestSecretValue,
			},
		},
	})

	p := HashicorpVaultProvider{
		Client: fakeClient,
	}
	res, err := p.GetSecret(fmt.Sprintf("%s:%s", TestSecretPath, TestSecretKey))
	assert.NoError(t, err)
	assert.Equal(t, TestSecretValue, res)

}

func TestHashicorpVaultProvider_GetSecretErr(t *testing.T) {
	p := HashicorpVaultProvider{
		Client: NewFakeVaultClient(),
	}
	_, err := p.GetSecret(fmt.Sprintf("%s:%s", TestSecretPath, TestSecretKey))
	assert.Error(t, err)
}

func TestHashicorpVaultProvider_GetKVSecrets(t *testing.T) {
	fakeClient := NewFakeVaultClient()

	fakeClient.SetSecret(TestSecretPath, &api.Secret{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				TestSecretKey: TestSecretValue,
			},
		},
	})

	p := HashicorpVaultProvider{
		Client: fakeClient,
	}

	res, err := p.GetKVSecrets(TestSecretPath)
	assert.NoError(t, err)
	assert.Equal(t,
		map[string]string{
			TestSecretKey: TestSecretValue,
		},
		res)
}

func TestHashicorpVaultProvider_GetKVSecretsErr(t *testing.T) {
	p := HashicorpVaultProvider{
		Client: NewFakeVaultClient(),
	}
	_, err := p.GetKVSecrets(TestSecretKey)
	assert.Error(t, err)
}
