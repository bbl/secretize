# Exec KRM Function Example

This example demonstrates how to use Secretize as an exec KRM function with Kustomize.

## Prerequisites

1. Install Kustomize (version 4.0.0 or later)
2. Build the secretize binary

## Setup

Build the secretize binary:

```bash
cd ../../..
go build -o secretize ./cmd/secretize
cd examples/exec/env
```

## Usage

1. Set the required environment variables:

```bash
export DATABASE_URL="postgresql://user:pass@localhost/db"
export API_KEY="your-secret-api-key"
export RENAMED_VAR="this-will-be-renamed"
export CONFIG_JSON='{"feature_new_ui": "true", "feature_beta": "false"}'
```

2. Run Kustomize build with KRM functions enabled:

```bash
kustomize build --enable-alpha-plugins --enable-exec .
```

## How it Works

The exec KRM function approach:

1. Kustomize recognizes the `config.kubernetes.io/function` annotation
2. It executes the specified binary path (`../../secretize`)
3. Kustomize sends a ResourceList to the binary's stdin
4. The binary processes the function config and returns modified ResourceList on stdout
5. The generated Secret is included in the final output

## Key Differences from Legacy

- Uses the KRM function specification
- Binary receives input via stdin/stdout instead of command-line arguments
- More flexible and follows the KRM standard
- Can process multiple resources in a pipeline

## Configuration

The `secret-generator.yaml` includes:
- `config.kubernetes.io/function`: Specifies the exec function path
- Standard SecretGenerator configuration for providers and secrets

## Example Output

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: exec-env-secrets
data:
  DATABASE_URL: cG9zdGdyZXNxbDovL3VzZXI6cGFzc0Bsb2NhbGhvc3QvZGI=
  API_KEY: eW91ci1zZWNyZXQtYXBpLWtleQ==
  newName: dGhpcy13aWxsLWJlLXJlbmFtZWQ=
  feature_flags: ZmFsc2U=
  new_ui: dHJ1ZQ==
  beta: ZmFsc2U=
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: exec-example-app
spec:
  # ... deployment spec ...
```

## Advanced Features

The exec KRM function supports:
- Renaming keys with the `newName=originalName` syntax
- Processing JSON values into multiple key-value pairs
- All provider types (env, aws-sm, azure-vault, hashicorp-vault, k8s-secret) 