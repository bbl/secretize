package generator

import (
	"encoding/base64"
	"strings"

	"github.com/DevOpsHiveHQ/secretize/internal/k8s"
	"github.com/DevOpsHiveHQ/secretize/internal/providers"
	"github.com/DevOpsHiveHQ/secretize/pkg/utils"
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kustomize/api/types"
)

type RegistryFunc func(params map[string]string) map[string]func() (providers.SecretsProvider, error)

func ProviderRegistry(params map[string]string) map[string]func() (providers.SecretsProvider, error) {
	return map[string]func() (providers.SecretsProvider, error){
		"aws-sm":          providers.NewAwsSMProvider,
		"hashicorp-vault": providers.NewHashicorpVaultProvider,
		"azure-vault": func() (providers.SecretsProvider, error) {
			return providers.NewAzureVaultProvider(params["name"])
		},
		"k8s-secret": func() (providers.SecretsProvider, error) {
			return providers.NewK8SSecretProvider(params["namespace"])
		},
		"env": func() (providers.SecretsProvider, error) {
			return providers.NewEnvProvider(), nil
		},
	}
}

type Literal struct {
	Key   string
	Value string
}

type SecretsSpec struct {
	KVLiterals []string  `yaml:"kv"`
	Literals   []Literal `yaml:"literals"`
}

type SecretSource struct {
	Provider    string `yaml:"provider"`
	SecretsSpec `yaml:",inline"`
	Params      map[string]string `yaml:"params"`
}

type SecretGenerator struct {
	Meta     types.ObjectMeta `yaml:"metadata"`
	Type     string           `yaml:"type"`
	Sources  []SecretSource   `yaml:"sources"`
	Literals []Literal        `yaml:"literals"`
}

func (l *Literal) UnmarshalYAML(unmarshal func(interface{}) error) error {
	stringLiteral := ""
	err := unmarshal(&stringLiteral)
	if err != nil {
		return err
	}
	l.Key = stringLiteral
	l.Value = stringLiteral

	if !strings.Contains(stringLiteral, "=") {
		return nil
	}

	// Split on the first "=" only
	idx := strings.Index(stringLiteral, "=")
	l.Key = stringLiteral[:idx]
	l.Value = stringLiteral[idx+1:]

	return nil
}

func ParseConfig(data []byte) (*SecretGenerator, error) {
	conf := SecretGenerator{}
	err := yaml.Unmarshal(data, &conf)
	return &conf, err
}

func (sg *SecretGenerator) FetchSecrets(registry RegistryFunc) (map[string]string, error) {
	secrets := make(map[string]string)

	for _, s := range sg.Sources {
		provider, err := registry(s.Params)[s.Provider]()
		if err != nil {
			return nil, err
		}
		providerSecrets, err := FetchProviderSecrets(provider, s.SecretsSpec)
		if err != nil {
			return nil, err
		}
		secrets = utils.Merge(secrets, providerSecrets)
	}
	return secrets, nil
}

func FetchProviderSecrets(p providers.SecretsProvider, spec SecretsSpec) (map[string]string, error) {
	res := make(map[string]string)

	for _, l := range spec.Literals {
		resp, err := p.GetSecret(l.Value)
		if err != nil {
			return nil, err
		}
		res[l.Key] = resp
	}

	for _, v := range spec.KVLiterals {
		resp, err := p.GetKVSecrets(v)
		if err != nil {
			return nil, err
		}
		res = utils.Merge(res, resp)
	}

	return res, nil
}

func (sg *SecretGenerator) Generate(secrets map[string]string) *k8s.Secret {
	secrets = utils.Map(secrets, func(v string) string {
		return base64.StdEncoding.EncodeToString([]byte(v))
	})
	return k8s.NewSecret(sg.Meta, sg.Type, secrets)
}
