package utils

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

func TestFatalErrCheck(t *testing.T) {
	exited := false
	logrus.StandardLogger().ExitFunc = func(i int) {
		exited = true
	}
	FatalErrCheck(errors.New("fatal error"))
	assert.True(t, exited)

	exited = false
	FatalErrCheck(nil)
	assert.False(t, exited)
}
