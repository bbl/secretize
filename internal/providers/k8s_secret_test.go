package providers

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
	v13 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	TestSecretName = "test"
)

func configureTestClient() (v13.SecretInterface, error) {
	k8sClient := testclient.NewSimpleClientset().CoreV1().Secrets("test")
	secretValue := make([]byte, base64.StdEncoding.EncodedLen(len([]byte(TestSecretValue))))
	base64.StdEncoding.Encode(secretValue, []byte(TestSecretValue))
	_, err := k8sClient.Create(context.Background(), &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: TestSecretName,
		},
		Data: map[string][]byte{
			TestSecretKey: secretValue,
		},
	}, metav1.CreateOptions{})
	return k8sClient, err
}

func TestK8SSecretProvider_GetSecret(t *testing.T) {
	k8sClient, err := configureTestClient()
	assert.NoError(t, err)
	p := K8SSecretProvider{
		Client: k8sClient,
	}

	res, err := p.GetSecret(fmt.Sprintf("%s:%s", TestSecretName, TestSecretKey))
	assert.NoError(t, err)
	assert.Equal(t, TestSecretValue, res)
}

func TestK8SSecretProvider_GetSecretErr(t *testing.T) {
	k8sClient, err := configureTestClient()
	assert.NoError(t, err)
	p := K8SSecretProvider{
		Client: k8sClient,
	}
	_, err = p.GetSecret("")
	assert.Error(t, err)
}

func TestK8SSecretProvider_GetKVSecrets(t *testing.T) {
	k8sClient, err := configureTestClient()
	assert.NoError(t, err)
	p := K8SSecretProvider{
		Client: k8sClient,
	}
	res, err := p.GetKVSecrets(TestSecretName)
	assert.NoError(t, err)
	assert.Equal(t,
		map[string]string{
			TestSecretKey: TestSecretValue,
		},
		res)
}

func TestK8SSecretProvider_GetKVSecretsErr(t *testing.T) {
	k8sClient, err := configureTestClient()
	assert.NoError(t, err)
	p := K8SSecretProvider{
		Client: k8sClient,
	}
	_, err = p.GetKVSecrets("")
	assert.Error(t, err)
}
