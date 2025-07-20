# Containerized KRM Function Example with HashiCorp Vault

This example demonstrates how to use Secretize as a **containerized KRM function** with Kustomize, configured to fetch secrets from HashiCorp Vault.

---

## Local Vault Testing with Docker Compose

A minimal `docker-compose.yml` is provided in this directory to spin up a Vault instance and pre-populate it with all the secrets needed for this example.

### Steps:
1. **Start Vault and initialize secrets:**
   ```bash
   docker-compose up -d
   # Wait for both 'vault' and 'setup' containers to finish initializing
   docker-compose ps
   # Vault UI: http://localhost:8200 (token: myroot)
   ```
2. **Build the Secretize Docker image (if using `image: secretize:local`):**
   ```bash
   cd ../../..
   docker build -t secretize:local .
   cd examples/docker/vault
   ```
   - If you use `image: ghcr.io/DevOpsHiveHQ/secretize:v0.1.0` in your function config, you can skip this step.
3. **Export the required Vault environment variables in your shell:**
   ```bash
   export VAULT_ADDR="http://127.0.0.1:8200"
   export VAULT_TOKEN="myroot"
   ```
   These variables will be used by the Secretize function to connect to Vault.
   - We use `127.0.0.1` since Vault is running on the host machine.
4. **Update the function config in `secret-generator.yaml` to set Vault environment variables:**
   ```yaml
   metadata:
     annotations:
       config.kubernetes.io/function: |
         container:
           image: secretize:local
           envs:
             - VAULT_ADDR
             - VAULT_TOKEN
   ```
   - We use `127.0.0.1` since Vault is running on the host machine.

5. **Run Kustomize build with containerized KRM function enabled:**
   ```bash
   kustomize build --enable-alpha-plugins --network .
   ```

---

## How it Works

- Kustomize recognizes the `config.kubernetes.io/function` annotation with `container` configuration.
- It pulls and runs the specified Docker image with Vault environment variables.
- The container connects to Vault using the provided token.
- Kustomize sends a ResourceList to the container's stdin.
- The container fetches secrets from Vault and returns the generated Secret on stdout.
- The Secret is included in the final output.

---

## Troubleshooting

### 1. **Secrets Not Found**
- Double-check the secret paths in `secret-generator.yaml`.
- If you have a subfolder called `data` in Vault, your path should be:
  ```yaml
  - DATABASE_URL=secret/data/data/docker-app/database-url:value
  ```
- If not, use:
  ```yaml
  - DATABASE_URL=secret/data/docker-app/database-url:value
  ```
- You can confirm the path with:
  ```bash
  vault kv get <full-path>
  # Example:
  vault kv get secret/data/docker-app/database-url
  vault kv get secret/data/data/docker-app/database-url
  ```
- If using the API, try:
  ```bash
  curl -H "X-Vault-Token: myroot" http://127.0.0.1:8200/v1/secret/data/docker-app/database-url
  curl -H "X-Vault-Token: myroot" http://127.0.0.1:8200/v1/secret/data/data/docker-app/database-url
  ```

### 2. **Authentication Failed**
- Make sure `VAULT_TOKEN` is set and valid in the function config.
- The default token for the test setup is `myroot`.

### 3. **Vault Not Reachable**
- Make sure `VAULT_ADDR` is set to `http://host.docker.internal:8200` (or `http://127.0.0.1:8200` if not using Docker Desktop).
- Ensure the Vault container is running and healthy.

### 4. **Image Not Found**
- If you see an error about `secretize:local` not found, make sure you built the image as described above.

---

## Technical Note: Vault Path Structure

- The full path to a secret is: `<mount>/data/<folder>/<secret>` for KV v2 API.
- In the Vault UI, you may see a subfolder called `data`—this is part of your logical path.
- If your secret is at `data/docker-app/api-key` in the UI under the `secret` mount, the full path is:
  ```yaml
  secret/data/docker-app/api-key
  ```
- If you have a subfolder called `data`, the path is:
  ```yaml
  secret/data/data/docker-app/api-key
  ```
- Always confirm with the Vault CLI or API if unsure.

---

## Security Considerations

- Never hardcode tokens in production.
- Use AppRole or Kubernetes authentication for production.
- Rotate tokens and audit secret access. 