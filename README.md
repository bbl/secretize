<p  align="center">
  <img width="409" height="450" src=".assets/logo.png"/>
</p>

Secretize is a kustomize plugin that helps to generate kubernetes secrets from various sources.  
It's like a swiss army knife, but for kubernetes secrets. 
 

---

## Sources

Secretize is able to generate secrets using the following providers:

- [AWS Secret Manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html)
- [Azure Vault](https://docs.microsoft.com/en-us/azure/key-vault/)
- [Hashicorp Vault](https://www.vaultproject.io/)
- [Other K8S secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- Environment variables

It is possible to use multiple providers at once.


## Installation

Install secretize to your `$XDG_CONFIG_HOME/kustomize/plugins` folder:

```bash
# todo
mkdir -p ...
curl github/releases... 
```



## Usage

There are two types of secrets: `literals` and `kv`. 
Literal secrets simply generate a single string output, while KV secrets will output with a dictionary of the key-value pairs. 
Some providers return only KV data by default (e.g. Hashicorp Vault and K8S Secret data), while others (AWS Secret Manager, Azure Vault and Environment variables) always output a single string. If you try to query `kv` secrets from later, secretize would treat the output as JSON and try to generate a `kv` output, e.g.:

```bash
export MY_KV_SECRET='{"secret_key_1":"secret_value_1", "secret_key_2": "secret_value_2"}'
export MY_LITERAL_SECRET=super_secret

cat <<'EOF' >./config.yaml
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: env-kv-secrets
sources:
    - provider: env
      kv:
        - MY_KV_SECRET
      literals: 
        - MY_LITERAL_SECRET
EOF
```

Secretize would generate the following secret using the config above:

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




### AWS Secret Manager

```yaml
apiVersion: secretize/v1
kind: SecretGenerator
metadata:
  name: my-aws-secrets
sources:
    - provider: aws-sm
      literals:
          - my-secret-1
          - SECRET_NAME=my-secret-1
```



