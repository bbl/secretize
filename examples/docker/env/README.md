# Containerized KRM Function Example (Environment Variables)

This example demonstrates how to use Secretize as a containerized KRM function with Kustomize, sourcing secrets from environment variables set directly in the function config.

---

## How It Works

- The containerized KRM function runs in an isolated environment and **does not inherit environment variables from your shell**.
- All required environment variables must be set explicitly in the `envs:` section of the function config in `secret-generator.yaml`.
- This approach is simple, reproducible, and works reliably with Kustomize and Secretize.

---

## Step-by-Step Usage

1. **(Optional) Build the Secretize Docker image (if using `image: secretize:local`):**
   ```bash
   cd ../../..
   docker build -t secretize:local .
   cd examples/docker/env
   ```
   
2. **Export the required environment variables in your shell:**
   ```bash
   export DATABASE_URL="postgresql://user:pass@localhost/db"
   export API_KEY="your-secret-api-key"
   export JWT_SECRET="your-jwt-secret"
   export CONFIG_JSON='{"feature_new_ui": "true", "feature_beta": "false"}'
   ```
   These variables will be referenced by the YAML configuration.

3. **Reference the environment variables in your YAML file:**
   In `secret-generator.yaml`, you can reference these variables using the `literals` and `kv` fields:
   ```yaml
   sources:
     - provider: env
       literals:
         - DATABASE_URL    # Reads from $DATABASE_URL
         - API_KEY         # Reads from $API_KEY
         - JWT_SECRET      # Reads from $JWT_SECRET
       kv:
         - CONFIG_JSON     # Reads from $CONFIG_JSON and parses as JSON
   ```
   - `literals`: Direct environment variable values
   - `kv`: Environment variables containing JSON that gets parsed into key-value pairs

4. **Run Kustomize build with containerized KRM function enabled:**
   ```bash
   kustomize build --enable-alpha-plugins .
   ```

---

## Why Not YAML Anchors for Env Vars?
- YAML anchors are useful for repeating static YAML blocks, but **they do not help with dynamic environment variable substitution**.
- For dynamic values, set them directly in the `envs:` list as shown above.
- If you want to use exported environment variables from your shell, use the [exec approach](../exec/env/README.md) instead.

---

## Troubleshooting

- If secrets are not found, make sure you set all required env vars in the function config.
- If you see an error about `secretize:local` not found, make sure you built the image as described above.
- If you see errors about `$` or `${VAR}` in the output, make sure you have replaced all placeholders with actual values.

---

## Security Considerations

- **Never hardcode real secrets in your configs for production.**
- Use this approach for local development, testing, or with non-sensitive values.
- For production, consider using a secrets manager (like Vault) and the appropriate Secretize provider.

---

## Reference: Using Host Environment Variables

- If you want to use environment variables exported in your shell, use the [exec KRM function approach](../exec/env/README.md), which can access your host environment directly. 