package providers

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

type AwsSMProvider struct {
	Client secretsmanageriface.SecretsManagerAPI
}

func (p *AwsSMProvider) GetSecret(name string) (string, error) {
	resp, err := p.Client.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &name})
	if err != nil {
		return "", err
	}
	return *resp.SecretString, err
}

func (p *AwsSMProvider) GetKVSecrets(name string) (map[string]string, error) {
	return jsonKVSecrets(p, name)
}

func NewAwsSMProvider() (SecretsProvider, error) {
	s, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return &AwsSMProvider{Client: secretsmanager.New(s)}, nil

}
