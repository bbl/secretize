package utils

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestMap(t *testing.T) {
	m := map[string]string{
		"k": "v",
	}
	res := Map(m, func(s string) string {
		return s
	})
	assert.Equal(t, res, m)
}

func TestMerge(t *testing.T) {
	m1 := map[string]string{
		"k": "v",
	}
	m2 := map[string]string{
		"k2": "v2",
	}
	res := Merge(m1, m2)
	assert.Equal(t, res, map[string]string{
		"k":  "v",
		"k2": "v2",
	})
}
