package main

import (
	"fmt"
	"io/ioutil"

	"github.com/bbl/secretize/pkg/generator"
	"github.com/bbl/secretize/pkg/utils"
	log "github.com/sirupsen/logrus"

	"os"
	"path/filepath"

	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func main() {
	// Check if running as KRM function (no args or stdin has content)
	if len(os.Args) == 1 || isKRMFunction() {
		runAsKRMFunction()
	} else {
		// Legacy mode
		runLegacyMode()
	}
}

// isKRMFunction checks if stdin has content (indicating KRM function mode)
func isKRMFunction() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// runLegacyMode runs the original secretize behavior
func runLegacyMode() {
	if len(os.Args) < 2 {
		log.Fatal("No argument passed, use `secretize /path/to/generator-config.yaml`")
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

// SecretGeneratorProcessor implements the KRM function processor
type SecretGeneratorProcessor struct{}

// Process implements the framework.ResourceListProcessor interface
func (p SecretGeneratorProcessor) Process(rl *framework.ResourceList) error {
	// Get the function config
	if rl.FunctionConfig == nil {
		return fmt.Errorf("no function config provided")
	}

	// Convert function config to YAML string
	fcString, err := rl.FunctionConfig.String()
	if err != nil {
		return fmt.Errorf("failed to marshal function config: %w", err)
	}

	// Parse as SecretGenerator
	secretGenerator, err := generator.ParseConfig([]byte(fcString))
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Generate secrets
	secrets, err := secretGenerator.FetchSecrets(generator.ProviderRegistry)
	if err != nil {
		return fmt.Errorf("failed to fetch secrets: %w", err)
	}

	// Generate the secret resource
	secret := secretGenerator.Generate(secrets)
	secretYaml, err := secret.ToYamlStr()
	if err != nil {
		return fmt.Errorf("failed to convert secret to yaml: %w", err)
	}

	// Parse the generated secret as RNode
	rNode, err := yaml.Parse(secretYaml)
	if err != nil {
		return fmt.Errorf("failed to parse generated secret: %w", err)
	}

	// Append to items
	rl.Items = append(rl.Items, rNode)

	return nil
}

// runAsKRMFunction runs secretize as a KRM function
func runAsKRMFunction() {
	processor := SecretGeneratorProcessor{}
	cmd := command.Build(processor, command.StandaloneDisabled, false)

	// Add dockerfile generation support
	command.AddGenerateDockerfile(cmd)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
