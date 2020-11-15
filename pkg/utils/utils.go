package utils

import log "github.com/sirupsen/logrus"

func Merge(m1, m2 map[string]string) map[string]string {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func Map(m map[string]string, f func(string) string) map[string]string {
	res := make(map[string]string)
	for k, v := range m {
		res[k] = f(v)
	}
	return res
}

func FatalErrCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
