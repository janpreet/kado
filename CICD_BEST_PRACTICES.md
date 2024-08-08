# Kado CI/CD Best Practices

This guide provides best practices for integrating Kado and Keybase into your CI/CD pipeline, especially when running inside containers.

## Table of Contents

1. [Mounting Directories](#mounting-directories)
2. [Keybase Integration](#keybase-integration)
3. [Security Considerations](#security-considerations)
4. [CI/CD Pipeline Configuration](#cicd-pipeline-configuration)
5. [Best Practices](#best-practices)
6. [Troubleshooting](#troubleshooting)

## Mounting Directories

To allow Kado to access your configuration files, templates, and other resources, you need to mount the relevant directories from your host system to the container.

### Docker Run Example

```bash
docker run -v $(pwd):/workspace ghcr.io/janpreet/kado:latest kado [command]
```

This command mounts the current directory to /kado-workspace in the container.

### Docker Compose Example

```yaml
yamlCopyversion: '3'
services:
  kado:
    image: ghcr.io/janpreet/kado:latest
    volumes:
      - ./:/workspace
    command: kado [command]
```

### CI/CD Configuration
In your CI/CD pipeline, ensure that your job checks out the repository and mounts it to the container:

```yaml
name: Deploy Infrastructure

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    container: 
      image: your-registry/kado:latest
      volumes:
        - ${{ github.workspace }}:/kado-workspace
    steps:
    - uses: actions/checkout@v2
    - name: Set up Keybase
      run: |
        echo "${{ secrets.KEYBASE_PAPERKEY }}" | keybase oneshot
        kado keybase link
      env:
        KEYBASE_PAPERKEY: ${{ secrets.KEYBASE_PAPERKEY }}
    - name: Deploy with Kado
      run: |
        cd /kado-workspace
        kado set cluster.yaml
```

## Keybase Integration

### Setting Up Keybase in CI/CD

1. Generate a paper key for your Keybase account.
2. Store the paper key securely in your CI/CD platform's secret management system.
3. In your CI/CD job, use the paper key to authenticate Keybase:

```yaml
job_name:
  image: your-registry/kado:latest
  script:
    - echo $KEYBASE_PAPERKEY | keybase oneshot
    - kado keybase link
    # Your Kado commands here
```

### Using Keybase Notes in Templates

Reference Keybase notes in your templates using the `{{keybase:note:note_name}}` syntax:

```hcl
pm_user = {{keybase:note:proxmox_api_key}}
pm_password = {{keybase:note:secret_token}}
```

## Security Considerations

1. **Paper Key Security**: Never expose your Keybase paper key in logs or non-secure storage.
2. **Ephemeral Sessions**: Use `keybase oneshot` to create temporary Keybase sessions.
3. **Least Privilege**: Use a Keybase account with minimal necessary permissions for CI/CD.
4. **Secure Note Storage**: Store sensitive information in Keybase notes, not in your codebase.

## CI/CD Pipeline Configuration

### Example GitHub Actions Workflow

```yaml
name: Deploy Infrastructure

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    container: your-registry/kado:latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Keybase
      run: |
        echo "${{ secrets.KEYBASE_PAPERKEY }}" | keybase oneshot
        kado keybase link
      env:
        KEYBASE_PAPERKEY: ${{ secrets.KEYBASE_PAPERKEY }}
    - name: Deploy with Kado
      run: kado set cluster.yaml
```

## Best Practices

1. **Configuration Management**:
   - Use `.kd` files for defining beads.
   - Keep sensitive data in Keybase notes, referenced in templates.

3. **Testing**:
   - Implement a test stage in your CI/CD pipeline using `kado -debug`.
   - Validate templates and configurations before deployment.

4. **Logging and Monitoring**:
   - Enable debug logging in CI/CD for troubleshooting.
   - Monitor Keybase activity for any suspicious actions.

5. **Secret Rotation**:
   - Regularly rotate your Keybase paper key.
   - Update Keybase notes with new credentials periodically.

6. **Error Handling**:
   - Implement proper error handling in your CI/CD scripts.
   - Set up notifications for pipeline failures.

## Troubleshooting

1. **Keybase Authentication Issues**:
   - Ensure the paper key is correctly stored in CI/CD secrets.
   - Check Keybase logs for authentication errors.

2. **Template Processing Errors**:
   - Verify that all referenced Keybase notes exist.
   - Check for syntax errors in your templates.

3. **Container Issues**:
   - Ensure all required tools are installed and accessible in the container.
   - Verify the container has necessary permissions to execute Kado and Keybase.

4. **Pipeline Failures**:
   - Review CI/CD logs for specific error messages.
   - Test Kado commands locally to replicate issues.

By following these best practices and guidelines, you can effectively integrate Kado and Keybase into your CI/CD pipeline, ensuring secure and efficient infrastructure management.
