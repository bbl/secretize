package k8s

import (
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/api/types"
	"testing"
)

const (
	Name = "mysecret"
	Namespace = "test"
	Type = "opaque"
)

func TestNewSecret(t *testing.T) {
	s := NewSecret(types.ObjectMeta{Namespace: Namespace, Name: Name}, Type, nil)
	assert.Equal(t, s.Type, Type)
	assert.Equal(t, s.Meta.Name, Name)
	assert.Equal(t, s.Meta.Namespace, Namespace)
}

func TestSecret_ToYamlStr(t *testing.T) {
	s := NewSecret(types.ObjectMeta{Namespace: Namespace, Name: Name}, Type, nil)
	_, err := s.ToYamlStr()
	assert.NoError(t, err)
}
