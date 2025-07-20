package generator

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kustomize/api/types"
)

func TestKRMIntegrationWithMultipleProviders(t *testing.T) {
	// Set up test environment variables
	os.Setenv("TEST_LITERAL", "literal-value")
	os.Setenv("TEST_JSON", `{"key1": "value1", "key2": "value2"}`)
	defer func() {
		os.Unsetenv("TEST_LITERAL")
		os.Unsetenv("TEST_JSON")
	}()

	// Test configuration with multiple sources
	config := &SecretGenerator{
		Meta: types.ObjectMeta{
			Name: "multi-source-secret",
		},
		Sources: []SecretSource{
			{
				Provider: "env",
				SecretsSpec: SecretsSpec{
					Literals: []Literal{
						{Key: "literal1", Value: "TEST_LITERAL"},
						{Key: "renamed", Value: "TEST_LITERAL"},
					},
					KVLiterals: []string{"TEST_JSON"},
				},
			},
		},
	}

	// Fetch secrets
	secrets, err := config.FetchSecrets(ProviderRegistry)
	assert.NoError(t, err)
	assert.NotEmpty(t, secrets)

	// Verify literals
	assert.Equal(t, "literal-value", secrets["literal1"])
	assert.Equal(t, "literal-value", secrets["renamed"])

	// Verify KV expansion
	assert.Equal(t, "value1", secrets["key1"])
	assert.Equal(t, "value2", secrets["key2"])

	// Generate the secret
	secret := config.Generate(secrets)
	assert.Equal(t, "multi-source-secret", secret.Meta.Name)

	// Verify base64 encoding
	for key, value := range secrets {
		encoded := base64.StdEncoding.EncodeToString([]byte(value))
		assert.Equal(t, encoded, secret.Data[key])
	}
}

func TestKRMFunctionConfigParsing(t *testing.T) {
	// Test parsing various valid configurations
	testCases := []struct {
		name   string
		yaml   string
		valid  bool
		errMsg string
	}{
		{
			name: "valid basic config",
			yaml: `
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: test
sources:
  - provider: env
    literals:
      - TEST_VAR
`,
			valid: true,
		},
		{
			name: "config with annotations",
			yaml: `
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: test
  annotations:
    config.kubernetes.io/function: |
      exec:
        path: ./secretize
sources:
  - provider: env
    literals:
      - TEST_VAR
`,
			valid: true,
		},
		{
			name: "empty sources",
			yaml: `
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: test
sources: []
`,
			valid: true,
		},
		{
			name: "multiple providers",
			yaml: `
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: test
sources:
  - provider: env
    literals:
      - VAR1
  - provider: env
    literals:
      - VAR2
`,
			valid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config, err := ParseConfig([]byte(tc.yaml))
			if tc.valid {
				assert.NoError(t, err)
				assert.NotNil(t, config)
			} else {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
			}
		})
	}
}

func TestLiteralParsing(t *testing.T) {
	testCases := []struct {
		input       string
		expectedKey string
		expectedVal string
	}{
		{"simple", "simple", "simple"},
		{"key=value", "key", "value"},
		{"key=value=with=equals", "key", "value=with=equals"},
		{"", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			var literal Literal
			err := yaml.Unmarshal([]byte(tc.input), &literal)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedKey, literal.Key)
			assert.Equal(t, tc.expectedVal, literal.Value)
		})
	}
}

func TestSecretTypeHandling(t *testing.T) {
	config := &SecretGenerator{
		Meta: types.ObjectMeta{
			Name: "typed-secret",
		},
		Type: "kubernetes.io/tls",
		Sources: []SecretSource{
			{
				Provider: "env",
				SecretsSpec: SecretsSpec{
					Literals: []Literal{
						{Key: "tls.crt", Value: "HOME"},
						{Key: "tls.key", Value: "USER"},
					},
				},
			},
		},
	}

	secrets, err := config.FetchSecrets(ProviderRegistry)
	assert.NoError(t, err)

	secret := config.Generate(secrets)
	assert.Equal(t, "kubernetes.io/tls", secret.Type)
}
