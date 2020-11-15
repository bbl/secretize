package providers

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewEnvProvider(t *testing.T) {
	p := NewEnvProvider()
	assert.NotNil(t, p)
}

func TestEnvProvider_GetSecret(t *testing.T) {
	p := NewEnvProvider()
	err := os.Setenv(TestKey, TestValue)
	assert.NoError(t, err)

	res, err := p.GetSecret(TestKey)
	assert.NoError(t, err)
	assert.Equal(t, TestValue, res)
}

func TestEnvProvider_GetSecretErr(t *testing.T) {
	p := NewEnvProvider()
	os.Unsetenv(TestKey)
	_, err := p.GetSecret(TestKey)
	assert.Error(t, err)
}

func TestEnvProvider_GetKVSecrets(t *testing.T) {
	p := NewEnvProvider()
	err := os.Setenv(TestKey, TestJsonStr)
	assert.NoError(t, err)

	res, err := p.GetKVSecrets(TestKey)
	assert.NoError(t, err)
	assert.Equal(t, TestSecretValue, res[TestSecretKey])
}

func TestEnvProvider_GetKVSecretsErr(t *testing.T) {
	p := NewEnvProvider()

	os.Unsetenv(TestKey)
	_, err := p.GetKVSecrets(TestKey)
	assert.Error(t, err)

	err = os.Setenv(TestKey, TestValue)
	assert.NoError(t, err)

	_, err = p.GetKVSecrets(TestKey)
	assert.Error(t, err)
}
