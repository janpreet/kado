# Kado Commands and Getting Started

## Available Commands

Kado supports several commands to help you manage and process your infrastructure configurations.

### `version`

Displays the current version of Kado.

```sh
kado version
```

### `config`

Loads and displays the bead configurations from the `*.kd` files in the current directory. This command shows the configuration with the order of execution.

```sh
kado config
```

### `set`

Processes the beads defined in the `*.kd` files and applies the configurations.

```sh
kado set
```

**Note:** If OPA (Open Policy Agent) is enabled and a bead is relayed to OPA for policy evaluation, you cannot run `kado set` without an approved policy. Beads that are not relayed to OPA or do not have policy enforcement can still be processed and set without OPA approval.

### `fmt`

Formats the `.kd` files in the proper Kado format. You can format all `.kd` files in the current directory or specify a single `.kd` file to format.

```sh
kado fmt
# or
kado fmt <filename.kd>
```

### `ai`

Analyzes the Terraform and Ansible configurations and provides infrastructure recommendations using an AI model. Requires AI configuration in `~/.kdconfig`.

```sh
kado ai
```

## Getting Started

### Installation

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/janpreet/kado.git
   cd kado
   ```

2. **Build the Binary**:
   ```sh
   make build
   ```

3. **Run Kado**:
   ```sh
   ./kado <command>
   ```

### Downloading Releases

You can download pre-built releases of Kado from the [GitHub Releases](https://github.com/janpreet/kado/releases) page.

### Configuration

Kado uses a configuration file located at `~/.kdconfig` for AI integration. The configuration file should include the following settings:

```sh
# ~/.kdconfig
AI_API_KEY=<your_api_key>
AI_MODEL=<your_model>
AI_CLIENT=<your_client_type>
AI_ENABLED=disabled # Can be overridden to enabled
```

- `AI_API_KEY`: Your API key for the AI model.
- `AI_MODEL`: The model you want to use (e.g., `gpt-3.5-turbo` or `claude-3-5-sonnet-20240620`).
- `AI_CLIENT`: The type of AI client (e.g., `chatgpt` or `anthropic_messages`).
- `AI_ENABLED`: Flag to enable or disable AI integration (default is `disabled`).

### Running Kado with AI

1. **Ensure AI is Enabled**:
   - Edit the `~/.kdconfig` file and set `AI_ENABLED=enabled`.

2. **Run the `ai` Command**:
   ```sh
   kado ai
   ```

   This command will analyze the Terraform and Ansible configurations and provide recommendations.

## Examples

### Example `.kd` File

```hcl
# Example Kado configuration
bead "example_bead" {
  key1 = "value1"
  key2 = "value2"
  #additional_comment
}
```

### Example AI Configuration

```sh
# ~/.kdconfig
AI_API_KEY=your_api_key
AI_MODEL=gpt-3.5-turbo
AI_CLIENT=chatgpt
AI_ENABLED=enabled
```

### Running Commands

- To format `.kd` files:
  ```sh
  kado fmt
  # or to format a specific file
  kado fmt example.kd
  ```

- To display bead configurations with the order of execution:
  ```sh
  kado config
  ```

- To process bead configurations and review IaC outputs:
  ```sh
  kado
  ```

- To apply IaC config:
  ```sh
  kado set
  ```

- To get AI recommendations:
  ```sh
  kado ai
  ```

## Command Usage Explanation

### Running `kado`

Running `kado` without any additional commands will start processing the bead configurations and apply the configurations found in the `*.kd` files in the current directory. For Ansible this will run playbook in dry-run mode, and Terraform is only until plan for user review.

### Running `kado set`

The `kado set` command processes the beads defined in the `*.kd` files and applies the configurations. This command ensures that the infrastructure setup is applied according to the definitions in the `.kd` files. This is equivalent of Terraform apply.

**Note:** If OPA (Open Policy Agent) is enabled and a bead is relayed to OPA for policy evaluation, you cannot run `kado set` without an approved policy. Beads that are not relayed to OPA or do not have policy enforcement can still be processed and set without OPA approval.

### Running `kado config`

The `kado config` command loads and displays the bead configurations from the `*.kd` files in the current directory. It shows the configuration with the order of execution. 
