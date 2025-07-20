<p  align="center">
  <img width="409" height="450" src=".assets/logo.png"/>
  <br>
<i> Secretize is a kustomize plugin that helps generating kubernetes secrets from various sources.  <br>
It's like a swiss army knife, but for kubernetes secrets. </i> 
  <br>
  <br>
  <img src="https://goreportcard.com/badge/github.com/DevOpsHiveHQ/secretize" />
<img src="https://github.com/DevOpsHiveHQ/secretize/workflows/CI/badge.svg">
   <a href="https://codecov.io/gh/DevOpsHiveHQ/secretize">
      <img src="https://codecov.io/gh/DevOpsHiveHQ/secretize/branch/main/graph/badge.svg" />
   </a>
  
</p>

---

## Sources

Secretize is able to generate secrets using the following providers:

- [AWS Secret Manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html)
- [Azure Vault](https://docs.microsoft.com/en-us/azure/key-vault/)
- [Hashicorp Vault](https://www.vaultproject.io/)
- [Other K8S secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [Environment variables](https://en.wikipedia.org/wiki/Environment_variable)

It is possible to use multiple providers at once.


## Installation

Secretize now supports multiple installation methods:

### Method 1: KRM Function (Recommended)

Secretize supports modern Kubernetes Resource Model (KRM) Functions, which work with Kustomize 4.0.0+:

#### Exec KRM Function
Download the binary and use it directly:
```bash
curl -L https://github.com/DevOpsHiveHQ/secretize/releases/download/v0.0.1/secretize-v0.0.1-linux-amd64.tar.gz | tar -xz
chmod +x secretize
```

#### Containerized KRM Function
Use the Docker image (no installation required):
```yaml
# In your kustomization, reference the container image
annotations:
  config.kubernetes.io/function: |
    container:
      image: ghcr.io/DevOpsHiveHQ/secretize:v0.1.0
```

### Method 2: Legacy Plugin (Deprecated)

Install secretize to your `$XDG_CONFIG_HOME/kustomize/plugin` folder:

1. Export the `XDG_CONFIG_HOME` variable if it's not already set:

```bash
export XDG_CONFIG_HOME=~/.config
```

2. Download the release binary into the kustomize plugin folder:

```bash
export SECRETIZE_DIR="$XDG_CONFIG_HOME/kustomize/plugin/secretize/v1/secretgenerator"
mkdir -p "$SECRETIZE_DIR"
curl -L https://github.com/DevOpsHiveHQ/secretize/releases/download/v0.0.1/secretize-v0.0.1-linux-amd64.tar.gz  | tar -xz -C $SECRETIZE_DIR
```

## Usage

### Using KRM Functions (Recommended)

With KRM functions, add the `config.kubernetes.io/function` annotation to your SecretGenerator:

#### Exec KRM Function Example
```yaml
# secret-generator.yaml
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: my-secrets
  annotations:
    config.kubernetes.io/function: |
      exec:
        path: ./secretize
sources:
  - provider: env
    literals:
      - DATABASE_URL
```

Run with: `kustomize build --enable-alpha-plugins --enable-exec .`

#### Containerized KRM Function Example
```yaml
# secret-generator.yaml
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: my-secrets
  annotations:
    config.kubernetes.io/function: |
      container:
        image: ghcr.io/DevOpsHiveHQ/secretize:v0.1.0 #TODO: Upload the image.
sources:
  - provider: env
    literals:
      - DATABASE_URL
```

Run with: `kustomize build --enable-alpha-plugins .`

### Legacy Plugin Usage

For the legacy plugin, use without annotations:

```yaml
# kustomization.yaml
generators:
  - secret-generator.yaml
```

Run with: `kustomize build --enable-alpha-plugins .`

### Provider Configuration

All providers can generate two types of secrets: `literals` and `kv` (Key-Value secrets).  
Literal secrets simply generate a single string output, while KV secrets will output with a dictionary of the key-value pairs.   

The full configuration API could be found in the [examples/secret-generator.yaml](./examples/secret-generator.yaml) file.

### AWS Secrets Manager

Fetching literal secrets is as simple, as using a default kustomize `secretGenerator` plugin:

```yaml
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: aws-sm-secrets
sources:
    - provider: aws-sm
      literals: 
        - mySecret
        - newName=mySecret 
```

The above config would query AWS Secrets Manager provider to get the `mySecret` string value. As a result, the following manifest will be generated:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: aws-sm-secrets
data:
  mySecret: c2VjcmV0X3ZhbHVlXzE= # a sample base64 encoded data 
  newName: c2VjcmV0X3ZhbHVlXzE=
```
 
Now let's assume that value of `mySecret` is a json string:
```json
{
  "secret_key_1":"secret_value_1", 
  "secret_key_2": "secret_value_2"
}
```

The generator config can be slightly modified, to generate a `kv` secret:

```yaml
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: aws-sm-secrets
sources:
    - provider: aws-sm
      kv: 
        - mySecret
```

As a result, the following secret is generated:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: aws-sm-secrets
data:
  secret_key_1: c2VjcmV0X3ZhbHVlXzE=
  secret_key_2: c2VjcmV0X3ZhbHVlXzI=
```

### Azure Vault

Azure Vault configuration is pretty similar to the above examples. However, there's additional `params` field, which is used to specify the Vault Name: 


```yaml
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: aws-sm-secrets
sources:
  - provider: azure-vault
    params:
      name: vault-name
    kv:
      - kv-secrets # will treat this as JSON, the same way as in the AWS example
    literals:
      - literal-secret-1
      - new_name=literal-secret-1
```


### Hashicorp Vault

Some providers only support key-value output, e.g. Hashicorp Vault and K8S Secret. 
For instance, the `mySecret` in Hashicorp Vault might look like the following:
```bash
vault kv get secret/mySecret
====== Data ======
Key           Value
---           -----
secret_key_1  secret_value_1
secret_key_2  secret_value_2
```

Querying provider's `kv` secrets will generate the corresponding key-value data:

```yaml
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: hashicorp-vault-secrets
sources:
    - provider: hashicorp-vault
      kv: 
        - secret/data/mySecret # you need to specify the full path in hashicorp vault provider
```
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: hashicorp-vault-secrets
data:
  secret_key_1: c2VjcmV0X3ZhbHVlXzE=
  secret_key_2: c2VjcmV0X3ZhbHVlXzI=
```

However you're able to query a certain literal in the key-value output using the following syntax: `secret-name:key`, e.g.:
  
```yaml
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: hashicorp-vault-secrets
sources:
    - provider: hashicorp-vault
      literals:
          - secret/data/mySecret-1:secret_key_1
```

As a result, the following manifest will be generated:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: hashicorp-vault-secrets
data:
  secret_key_1: c2VjcmV0X3ZhbHVlXzE=
```

## Kubernetes Secret

Kubernetes secret provider is similar to the Hashicorp Vault. Additionally, this provider expects the `params` field with the `namespace` specification.   
You're able to get the entire secret data using the `kv` query, or get a particular key using the `literals` query with the `:` delimiter syntax:

```yaml
# The original secret in a default namespace
#
apiVersion: v1
kind: Secret
metadata:
  name: original-secret
  namespace: default
data:
  secret_key_1: c2VjcmV0X3ZhbHVlXzE=
  secret_key_2: c2VjcmV0X3ZhbHVlXzI=
---
# Secret generator configuration
#
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: kubernetes-secrets
sources:
    - provider: k8s-secret
      params:
        namespace: default
      kv:
        - original-secret
      literals:
        - new_name=original-secret:secret_key_1
---
# Generated secret
#
apiVersion: v1
kind: Secret
metadata:
  name: kubernetes-secrets
data:
  secret_key_1: c2VjcmV0X3ZhbHVlXzE=
  secret_key_2: c2VjcmV0X3ZhbHVlXzI=
  new_name: c2VjcmV0X3ZhbHVlXzE=

```
 

## Env 

The environment variables plugin is similar to the AWS and Azure plugins. The `literals` would simply fetch corresponding environment variables, while `kv` would treat each variable as JSON and try to parse it:

```yaml
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: env-secrets
sources:
    - provider: env
      kv:
        - MY_KV_SECRET
      literals: 
        - MY_LITERAL_SECRET
```

Secretize will fetch the corresponding environment variables during the `kustomize build` command:

```bash
export MY_KV_SECRET='{"secret_key_1":"secret_value_1", "secret_key_2": "secret_value_2"}'
export MY_LITERAL_SECRET=super_secret

kustomize build
```

The following secret is generated:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: env-kv-secrets
data:
  MY_LITERAL_SECRET: c3VwZXJfc2VjcmV0
  secret_key_1: c2VjcmV0X3ZhbHVlXzE=
  secret_key_2: c2VjcmV0X3ZhbHVlXzI=
```

## Examples

Check out the [examples](./examples) directory for complete working examples:

- [Legacy Plugin Example](./examples/legacy) - Traditional Kustomize plugin approach
- [Exec KRM Function Example](./examples/exec) - Modern exec-based KRM function
- [Containerized KRM Function Example](./examples/docker) - Docker-based KRM function

## Test Infrastructure

For comprehensive testing with real secret stores, see the [test-infrastructure](./test-infrastructure/) directory which provides:

- **HashiCorp Vault** setup with test secrets
- **AWS Secrets Manager** emulation via LocalStack  
- **Kubernetes** cluster with test secrets
- **Automated testing** for all providers and execution modes

```bash
cd test-infrastructure
./test-all-providers.sh
```

## Documentation

For detailed documentation on KRM Functions support, see [KRM Functions Documentation](./docs/KRM_FUNCTIONS.md).
