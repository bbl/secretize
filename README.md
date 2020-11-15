<p  align="center">
  <img width="409" height="450" src=".assets/logo.png"/>
<i> Secretize is a kustomize plugin that helps to generate kubernetes secrets from various sources.  
It's like a swiss army knife, but for kubernetes secrets. </i>
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

Install secretize to your `$XDG_CONFIG_HOME/kustomize/plugins` folder:

```bash
# todo
mkdir -p ...
curl github/releases... 
```



## Usage

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
