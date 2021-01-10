package providers

import (
	"fmt"
	"github.com/hashicorp/vault/api"
)

type VaultReader interface {
	Read(string) (*api.Secret, error)
}

type HashicorpVaultProvider struct {
	Client VaultReader
}

func (p *HashicorpVaultProvider) GetKVSecrets(path string) (map[string]string, error) {
	s, err := getHashicorpVaultSecret(p.Client, path)
	if err != nil {
		return nil, err
	}
	kvSecrets := make(map[string]string)
	for k, v := range s.Data["data"].(map[string]interface{}) {
		kvSecrets[k] = v.(string)
	}
	return kvSecrets, err
}

func NewHashicorpVaultProvider() (SecretsProvider, error) {
	vaultToken, err := getToken()
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	client.SetToken(vaultToken)

	return &HashicorpVaultProvider{Client: client.Logical()}, err
}

func getToken() (string, error) {
	t, exists := os.LookupEnv("VAULT_TOKEN")
	if !exists {
		tokenPath, exists := os.LookupEnv("VAULT_TOKEN_PATH")
		if exists == false {
			return "", errors.New("$VAULT_TOKEN unset and $VAULT_TOKEN_PATH unset")
		}

		tBytes, err := ioutil.ReadFile(tokenPath)
		if err != nil {
			fmt.Println("$VAULT_TOKEN_PATH set, but failed to read file")
			return "", err
		}

		t = strings.TrimSpace(string(tBytes))
		if len(t) == 0 {
			fmt.Println("$VAULT_TOKEN_PATH set, but file is empty")
			return "", err
		}
	}

	return t, nil
}

func getKVSecret(name string, f func(path string) (map[string]string, error)) (string, error) {
	path, field := parseSecretLiteral(name)
	if field == "" {
		return "", fmt.Errorf("vault field is empty")
	}
	data, err := f(path)
	if err != nil {
		return "", err
	}
	if val, exists := data[field]; exists {
		return val, err
	}
	return "", fmt.Errorf("field \"%s\" not present in secret %s", field, path)
}

func (p *HashicorpVaultProvider) GetSecret(name string) (string, error) {
	return getKVSecret(name, func(path string) (map[string]string, error) {
		s, err := getHashicorpVaultSecret(p.Client, path)
		if err != nil {
			return nil, err
		}
		kvSecrets := make(map[string]string)
		for k, v := range s.Data["data"].(map[string]interface{}) {
			kvSecrets[k] = v.(string)
		}
		return kvSecrets, nil
	})
}

func getHashicorpVaultSecret(client VaultReader, path string) (*api.Secret, error) {
	s, err := client.Read(path)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, fmt.Errorf("couldn't find the spicified secret mount: %s", path)
	}
	if s.Data == nil {
		return nil, fmt.Errorf("no value found at %s", path)
	}
	return s, err
}
