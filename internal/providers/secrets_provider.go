package providers

import (
	"encoding/json"
	"strings"
)

const (
	secretNameSeparator = ":"
)

func parseSecretLiteral(literal string) (path string, field string) {
	res := strings.Split(literal, secretNameSeparator)
	path = res[0]
	if len(res) > 1 {
		field = res[1]
	}
	return path, field
}

func jsonKVSecrets(p SecretsProvider, name string) (map[string]string, error) {
	secret, err := p.GetSecret(name)
	if err != nil {
		return nil, err
	}
	kvSecrets := make(map[string]string)
	err = json.Unmarshal([]byte(secret), &kvSecrets)
	return kvSecrets, err
}

type SecretsProvider interface {
	GetSecret(name string) (string, error)
	GetKVSecrets(path string) (map[string]string, error)
}
