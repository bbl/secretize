package generator

import (
	"encoding/base64"
	"fmt"
	"github.com/bbl/kustomize-secrets/internal/providers"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/api/types"
	"testing"
)

var (
	testSG = SecretGenerator{
		Meta:     types.ObjectMeta{},
		Type:     "",
		Sources:  nil,
		Literals: nil,
	}
	testSecretsSpec = SecretsSpec{
		Literals: []Literal{
			{Key: "secret-key-name", Value: "provider-key-name"},
		},
		KVLiterals: []string{
			"test-kv",
		}}
	testProvider = &fakeProvider{
		literals: map[string]string{
			"provider-key-name": "test-value",
		},
		kv: map[string]map[string]string{
			"test-kv": {},
		},
	}
)

func TestParseConfig(t *testing.T) {
	confStr := `metadata:
  name: infrabot-secrets
  annotations:
    kustomize.config.k8s.io/behavior: "merge"
sources:
  - provider: hashicorp-vault
    literals:
      - NAME=example:example
`
	conf, err := ParseConfig([]byte(confStr))
	assert.NoError(t, err)
	assert.Equal(t, "hashicorp-vault", conf.Sources[0].Provider)
}

type fakeProvider struct {
	literals map[string]string
	kv       map[string]map[string]string
}

func (f *fakeProvider) GetSecret(name string) (string, error) {
	if val, ok := f.literals[name]; ok {
		return val, nil
	}
	return "", fmt.Errorf("literal key not found error: %s", name)
}
func (f *fakeProvider) GetKVSecrets(path string) (map[string]string, error) {
	if val, ok := f.kv[path]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("kv path not found error, %s", path)
}

func TestFetchProviderSecrets(t *testing.T) {
	res, err := FetchProviderSecrets(testProvider, testSecretsSpec)

	assert.NoError(t, err)
	assert.Equal(t, "test-value", res["secret-key-name"])

}

func TestFetchProviderSecretsErr(t *testing.T) {
	_, err := FetchProviderSecrets(&fakeProvider{}, testSecretsSpec)
	assert.Error(t, err)

	_, err = FetchProviderSecrets(&fakeProvider{
		literals: map[string]string{
			"provider-key-name": "test-value",
		},
	}, testSecretsSpec)
	assert.Error(t, err)

}

func TestSecretGenerator_FetchSecrets(t *testing.T) {
	testSG.Sources = append(testSG.Sources, SecretSource{
		Provider:    "fake",
		SecretsSpec: testSecretsSpec,
		Params:      nil,
	})
	_, err := testSG.FetchSecrets(func(params map[string]string) map[string]func() (providers.SecretsProvider, error) {
		return map[string]func() (providers.SecretsProvider, error){
			"fake": func() (providers.SecretsProvider, error) {
				return testProvider, nil
			},
		}
	})
	assert.NoError(t, err)
}

func TestSecretGenerator_FetchSecretsErr(t *testing.T) {
	testSG.Sources = append(testSG.Sources, SecretSource{
		Provider:    "fake",
		SecretsSpec: testSecretsSpec,
		Params:      nil,
	})

	_, err := testSG.FetchSecrets(func(params map[string]string) map[string]func() (providers.SecretsProvider, error) {
		return map[string]func() (providers.SecretsProvider, error){
			"fake": func() (providers.SecretsProvider, error) {
				return nil, fmt.Errorf("error")
			},
		}
	})
	assert.Error(t, err)
	assert.Equal(t, "error", err.Error())

	_, err = testSG.FetchSecrets(func(params map[string]string) map[string]func() (providers.SecretsProvider, error) {
		return map[string]func() (providers.SecretsProvider, error){
			"fake": func() (providers.SecretsProvider, error) {
				return &fakeProvider{}, nil
			},
		}
	})
	assert.Error(t, err)
}

func TestProviderRegistry(t *testing.T) {
	registry := ProviderRegistry(map[string]string{})
	assert.NotEmpty(t, registry)
}

func TestSecretGenerator_Generate(t *testing.T) {
	s := testSG.Generate(map[string]string{
		"test-key": "test-value",
	})
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("test-value")), s.Data["test-key"])
}
