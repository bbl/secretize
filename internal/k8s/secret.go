package k8s

import (
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kustomize/api/types"
)

const (
	ApiV1      = "v1"
	SecretKind = "Secret"
)

type Secret struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Meta       types.ObjectMeta  `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
	Type       string            `yaml:"type,omitempty"`
}

func NewSecret(meta types.ObjectMeta, secretType string, data map[string]string) *Secret {
	return &Secret{
		ApiVersion: ApiV1,
		Kind:       SecretKind,
		Meta:       meta,
		Type:       secretType,
		Data:       data,
	}
}

func (s *Secret) ToYamlStr() (string, error) {
	out, err := yaml.Marshal(s)
	return string(out), err
}
