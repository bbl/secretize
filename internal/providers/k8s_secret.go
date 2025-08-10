package providers

import (
	"context"
	"encoding/base64"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type K8SSecretProvider struct {
	Client v1.SecretInterface
}

func Base64Decode(src []byte) ([]byte, error) {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	_, err := base64.StdEncoding.Decode(decoded, src)
	return decoded, err
}

func (p *K8SSecretProvider) GetSecret(name string) (string, error) {
	return getKVSecret(name, func(path string) (map[string]string, error) {
		s, err := p.Client.Get(context.Background(), path, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		res := make(map[string]string)
		for k, v := range s.Data {
			decoded, err := Base64Decode(v)
			if err != nil {
				return nil, err
			}
			res[k] = string(decoded)
		}
		return res, nil

	})
}

func (p *K8SSecretProvider) GetKVSecrets(path string) (map[string]string, error) {
	res, err := p.Client.Get(context.Background(), path, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	kvSecrets := make(map[string]string)
	for k, v := range res.Data {
		decoded, err := Base64Decode(v)
		if err != nil {
			return nil, err
		}
		kvSecrets[k] = string(decoded)
	}
	return kvSecrets, err
}

func NewK8SSecretProvider(namespace string) (SecretsProvider, error) {
	kubeConfigPath := os.Getenv("KUBECONFIG")

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	v1Client, err := v1.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8SSecretProvider{Client: v1Client.Secrets(namespace)}, nil
}
