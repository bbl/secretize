package providers

import (
	"fmt"
	"os"
)

type EnvProvider struct {
}

func (p *EnvProvider) GetSecret(name string) (string, error) {
	if res, found := os.LookupEnv(name); found {
		return res, nil
	}
	return "", fmt.Errorf("couldn't find env variable: %s", name)
}

func (p *EnvProvider) GetKVSecrets(name string) (map[string]string, error) {
	return jsonKVSecrets(p, name)
}

func NewEnvProvider() SecretsProvider {
	return &EnvProvider{}

}
