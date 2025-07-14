# Legacy Kustomize Plugin Example with HashiCorp Vault

⚠️ **DEPRECATED** ⚠️

This example demonstrates the **legacy** Kustomize plugin system with HashiCorp Vault, which is **deprecated** and may not work with newer versions of Kustomize (v4.0.0+).

## ⚠️ Important Notice

**This approach is deprecated and may not work with current Kustomize versions.** The legacy plugin system has been replaced by the Kubernetes Resource Model (KRM) Functions.

### Recommended Alternative

For new projects, use the **KRM Function approach** instead:
- See [`../../exec/vault/`](../../exec/vault/) for a working example

The modern approach is more reliable, follows standards, and works with current Kustomize versions.

## Why This May Not Work

The legacy plugin system has several issues with newer Kustomize versions:

1. **Plugin Interface Changes**: The plugin execution interface has changed, causing "no function config provided" errors
2. **Deprecated Architecture**: The `$XDG_CONFIG_HOME/kustomize/plugin/` directory structure is no longer the recommended approach
3. **Limited Compatibility**: May not work with Kustomize v4.0.0+ due to architectural changes

### Technical Root Cause

The fundamental issue is that **Kustomize v4.0.0+ changed how it calls plugins**:
- **Older Kustomize versions (v3.x)**: Called plugins with command-line arguments
- **Newer Kustomize versions (v4.0.0+)**: Call plugins via stdin in KRM Function format

The Secretize plugin detects stdin input and switches to KRM Function mode, but then fails because:
- Kustomize sends the config via stdin
- The plugin expects a `ResourceList` format with `FunctionConfig` field
- But Kustomize is sending just the raw YAML config
- The plugin fails with "no function config provided" because it can't find the `FunctionConfig` field

**Why direct execution works**: When calling `"$SECRETIZE_DIR/SecretGenerator" secret-generator.yaml` directly, it uses the legacy mode (command-line arguments), which still works perfectly.

This explains why the legacy plugin system is fundamentally incompatible with newer Kustomize versions - the interface has changed completely.

---

## Prerequisites

1. Install Kustomize (version 3.x or earlier for best compatibility)
2. Install Docker
3. Install Secretize plugin to the Kustomize plugin directory

---

## Local Vault Testing with Minimal Docker Compose

You can use the provided `docker-compose.yml` in this folder to spin up a local Vault instance with all the secrets needed for this example.

### Steps:

1. **Start Vault and initialize secrets:**
   ```bash
   docker-compose up -d
   # Wait for both 'vault' and 'setup' containers to finish initializing
   docker-compose ps
   # Vault UI: http://localhost:8200 (token: myroot)
   ```

2. **Set up the legacy plugin:**
   ```bash
   export XDG_CONFIG_HOME=~/.config
   export SECRETIZE_DIR="$XDG_CONFIG_HOME/kustomize/plugin/secretize/v1/secretgenerator"
   mkdir -p "$SECRETIZE_DIR"
   go build -o "$SECRETIZE_DIR/SecretGenerator" ../../../cmd/secretize
   ```

3. **Run the plugin directly (recommended for legacy):**
   
   Before running, set the Vault address and token environment variables to match your local Vault instance:
   ```bash
   export VAULT_ADDR="http://127.0.0.1:8200"
   export VAULT_TOKEN="myroot"
   "$SECRETIZE_DIR/SecretGenerator" secret-generator.yaml
   ```
   This will output the generated Kubernetes Secret using the secrets from Vault.

4. **(Optional) Try Kustomize build (may fail on modern versions):**
   ```bash
   kustomize build --enable-alpha-plugins .
   ```
   **Note:** This will likely fail with modern Kustomize due to the reasons explained above.

---

## Vault Secrets Used

The setup container in `docker-compose.yml` will create these secrets:
- `secret/data/docker-app/database-url`
- `secret/data/docker-app/api-key`
- `secret/data/docker-app/jwt-secret`
- `secret/data/docker-app/app-config`
- `secret/data/docker-app/feature-flags`

---

## Example Output

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: legacy-vault-secrets
data:
  DATABASE_URL: <base64>
  API_KEY: <base64>
  JWT_SECRET: <base64>
  ...
```

---

## Migration to Modern Approach

To migrate from this legacy approach to the modern KRM Function approach, see [`../../exec/vault/`](../../exec/vault/). 