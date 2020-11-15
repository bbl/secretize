package main

import (
	"fmt"
	"github.com/bbl/kustomize-secrets/pkg/generator"
	"github.com/bbl/kustomize-secrets/pkg/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"

	"os"
	"path/filepath"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal(
			"No argument passed, use `secretize /path/to/generator-config.yaml`")
	}

	filename, _ := filepath.Abs(os.Args[1])
	yamlFile, err := ioutil.ReadFile(filename)
	utils.FatalErrCheck(err)

	secretGenerator, err := generator.ParseConfig(yamlFile)
	utils.FatalErrCheck(err)

	secrets, err := secretGenerator.FetchSecrets(generator.ProviderRegistry)
	utils.FatalErrCheck(err)

	s := secretGenerator.Generate(secrets)
	out, err := s.ToYamlStr()
	utils.FatalErrCheck(err)
	fmt.Println(out)
}
