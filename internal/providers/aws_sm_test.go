package providers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeSMClient struct {
	secretsmanageriface.SecretsManagerAPI
	data map[string]string
}

var (
	EmptySecretError = fmt.Errorf("secret value is missing")
	TestKey          = "key"
	TestValue        = "value"
	TestSecretKey    = "secret-key"
	TestSecretValue  = "secret-value"
	TestJsonStr      = fmt.Sprintf(`{"%s":"%s"}`, TestSecretKey, TestSecretValue)
	TestMap          = map[string]string{
		TestKey: TestValue,
	}
	TestMapJson = map[string]string{
		TestKey: TestJsonStr,
	}
)

func (m *fakeSMClient) setData(data map[string]string) {
	m.data = data
}

func newFakeSMClient(data map[string]string) *fakeSMClient {
	return &fakeSMClient{data: data}
}

func (m *fakeSMClient) GetSecretValue(in *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	if val, found := m.data[*in.SecretId]; found {
		return &secretsmanager.GetSecretValueOutput{SecretString: &val}, nil

	}
	return nil, EmptySecretError
}

func TestAwsSMProvider_GetSecret(t *testing.T) {

	p := AwsSMProvider{
		Client: newFakeSMClient(TestMap),
	}
	res, err := p.GetSecret(TestKey)
	assert.NoError(t, err)
	assert.Equal(t, TestValue, res)
}

func TestAwsSMProvider_GetSecretErr(t *testing.T) {
	p := AwsSMProvider{
		Client: newFakeSMClient(nil),
	}
	_, err := p.GetSecret(TestKey)
	assert.Error(t, err)
}

func TestAwsSMProvider_GetKVSecrets(t *testing.T) {
	p := AwsSMProvider{
		Client: newFakeSMClient(TestMapJson),
	}
	res, err := p.GetKVSecrets(TestKey)
	assert.NoError(t, err)
	assert.Equal(t, TestSecretValue, res[TestSecretKey])
}

func TestAwsSMProvider_GetKVSecretsErr(t *testing.T) {
	mockClient := newFakeSMClient(map[string]string{})
	p := AwsSMProvider{
		Client: mockClient,
	}
	_, err := p.GetKVSecrets(TestKey)
	assert.Error(t, err)

	// non-json value
	mockClient.setData(TestMap)
	_, err = p.GetKVSecrets(TestKey)
	assert.Error(t, err)

}

func TestNewAwsSMProvider(t *testing.T) {
	_, err := NewAwsSMProvider()
	assert.NoError(t, err)
}
