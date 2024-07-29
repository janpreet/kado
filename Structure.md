Here is the detailed documentation for the codebase based on the provided structure:

## Project Structure

```plaintext
.
├── Configuration.md
├── Dockerfile
├── How-to.md
├── LandingZone
│   └── Empty.txt
├── Makefile
├── README.md
├── Structure.md
├── VERSION
├── bump_version.py
├── examples
│   ├── cluster.kd
│   ├── cluster.yaml
│   ├── relay.kd
│   └── templates
│       ├── ansible
│       │   └── inventory.tmpl
│       └── terraform
│           ├── backend.tfvars.tmpl
│           └── vm.tfvars.tmpl
├── go.mod
├── go.sum
├── kado
├── main.go
├── packages
│   ├── ansible
│   │   └── ansible.go
│   ├── bead
│   │   └── bead.go
│   ├── config
│   │   ├── beadconfig.go
│   │   ├── config.go
│   │   └── yamlconfig.go
│   ├── display
│   │   └── display.go
│   ├── engine
│   │   ├── ai.go
│   │   ├── engine.go
│   │   └── formatter.go
│   ├── helper
│   │   ├── gitclone.go
│   │   └── helper.go
│   ├── opa
│   │   └── opa.go
│   ├── render
│   │   ├── driver.go
│   │   ├── kd.go
│   │   ├── writer.go
│   │   └── yaml.go
│   └── terraform
│       └── terraform.go
├── templates
│   ├── ansible
│   │   └── inventory.tmpl
│   └── terraform
│       ├── backend.tfvars.tmpl
│       └── terraform.tfvars.tmpl
```

## File and Directory Overview

### Root Directory

- **LandingZone/**: Directory where the repositories and files required by the beads are cloned and processed.
- **cluster.kd**: The main custom configuration file that defines the beads and their properties.
- **cluster.yaml**: YAML configuration file for the cluster setup.
- **go.mod** and **go.sum**: Go modules files for dependency management.
- **main.go**: The main entry point of the application.
- **readme.md**: Documentation file for the project.
- **testing.kd**: An additional KD file for testing purposes.

### Packages Directory

#### Ansible

- **ansible.go**: Contains functions to handle the execution of Ansible playbooks.

#### Bead

- **bead.go**: Defines the structure and properties of a bead.

#### Config

- **beadconfig.go**: Contains functions to load bead configurations.
- **config.go**: Contains functions to load general configurations.
- **yamlconfig.go**: Contains functions to load and parse YAML configurations.

#### Display

- **display.go**: Contains functions to display bead configurations and YAML content.

#### Engine

- **engine.go**: Contains the main function for handling the execution of Ansible playbooks.

#### Helper

- **gitclone.go**: Contains functions to clone Git repositories.
- **helper.go**: Contains helper functions for various operations such as file checks and setting up the environment.

#### OPA

- **opa.go**: Contains functions to handle OPA policy evaluation and related actions.

#### Render

- **driver.go**: Contains the main driver functions for rendering templates.
- **kd.go**: Contains functions to parse and process KD files.
- **writer.go**: Contains functions to write output files.
- **yaml.go**: Contains functions to handle YAML processing.

#### Terraform

- **terraform.go**: Contains functions to handle Terraform operations such as planning and applying configurations.

### Templates Directory

#### Ansible

- **inventory.tmpl**: Template for the Ansible inventory.

#### Terraform

- **backend.tfvars.tmpl**: Template for Terraform backend variables.
- **vm.tfvars.tmpl**: Template for Terraform VM variables.

## Detailed Documentation

### main.go

This is the main entry point of the application. It orchestrates the processing of beads, loading configurations, and setting up the environment.

Key Functions:

- **main**: Entry point of the application. Handles command-line arguments and processes beads.
- **processBead**: Processes a single bead, including cloning repositories, rendering templates, and handling Ansible and Terraform operations.
- **convertYAMLToSlice**: Converts YAML data to a slice of maps.
- **applyRelayOverrides**: Applies overrides for relay fields.

### packages/bead/bead.go

Defines the structure and properties of a bead. A bead represents a unit of work or configuration in the system.

### packages/config/config.go

Contains functions to load general configurations, including bead configurations and YAML configurations.

Key Functions:

- **LoadBeadsConfig**: Loads bead configurations from a custom format file.
- **LoadYAMLConfig**: Loads and parses YAML configuration files.

### packages/display/display.go

Contains functions to display bead configurations and YAML content.

Key Functions:

- **DisplayBeads**: Displays parsed beads from KD files.
- **DisplayYAMLs**: Displays parsed YAML content.
- **DisplayTemplateOutput**: Displays the result of processing templates.
- **DisplayBeadConfig**: Displays the configuration and order of execution of beads.

### packages/engine/engine.go

Contains the main function for handling the execution of Ansible playbooks.

Key Functions:

- **HandleAnsible**: Executes an Ansible playbook with the given configuration.

### packages/helper/helper.go

Contains helper functions for various operations such as file checks and setting up the environment.

Key Functions:

- **FileExists**: Checks if a file exists.
- **SetupLandingZone**: Sets up the LandingZone directory.
- **CloneRepo**: Clones a Git repository.

### packages/opa/opa.go

Contains functions to handle OPA (Open Policy Agent) policy evaluation and related actions.

Key Functions:

- **HandleOPA**: Handles the processing of the OPA bead, including policy evaluation and action handling.

### packages/render/kd.go

Contains functions to parse and process KD files.

Key Functions:

- **GetKDFiles**: Gets all KD files in the specified directory.
- **ProcessKdFiles**: Processes all KD files and returns beads and invalid bead names.
- **parseKdFile**: Parses a single KD file and returns beads and invalid bead names.

### packages/terraform/terraform.go

Contains functions to handle Terraform operations such as planning and applying configurations.

Key Functions:

- **HandleTerraform**: Handles the processing of the Terraform bead, including planning and applying configurations.

## Bead Structure and Valid Fields

### General Structure

A bead is defined in a KD file with the following structure:

```plaintext
bead "<bead_name>" {
  <field1> = "<value1>"
  <field2> = "<value2>"
  ...
}
```

### Valid Fields for Specific Beads

#### Ansible Bead

- **source**: URL of the Git repository containing the Ansible playbook.
- **playbook**: Path to the Ansible playbook file.
- **extra_vars_file**: Boolean indicating if extra variables file should be used.
- **relay**: Name of the relay bead.
- **relay_field**: Overrides for relay fields.

#### Terraform Bead

- **source**: URL of the Git repository containing the Terraform configuration.
- **relay**: Name of the relay bead.
- **relay_field**: Overrides for relay fields.

#### OPA Bead

- **path**: Path to the OPA policy file.
- **input**: Path to the input file for OPA evaluation.
- **package**: OPA package to evaluate.

#### Example Bead Definition

```plaintext
bead "ansible" {
  enabled = true
  source = "git@github.com:janpreet/proxmox_ansible.git"
  playbook = "cluster.yaml"
  extra_vars_file = false
  relay = opa
  relay_field = "source=git@github.com:janpreet/proxmox_ansible.git,path=ansible/policies/proxmox.rego,input=ansible/cluster.yaml,package=data.proxmox.main.allowed"
}
```

## Additional Notes

- **Relay Mechanism**: Beads can relay execution to another bead using the `relay` and `relay_field` properties.
- **OPA Evaluation**: OPA policies are evaluated to determine if actions (e.g., Terraform apply, Ansible playbook execution) should be allowed or denied.
- **Command-line Arguments**: The application can be run with different command-line arguments (e.g., `version`, `config`, `set`) to control its behavior.

This detailed documentation should help in understanding the structure, functionality, and usage of the codebase. If you have any further questions or need additional information, feel free to ask!