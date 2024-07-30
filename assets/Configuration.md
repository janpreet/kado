# Kado Configuration Documentation

## Bring Your Own Code (BYOC)

Kado is designed to be flexible and adaptable to your infrastructure needs. Users can plug the source for their parameterized Infrastructure as Code (IaC) into Kado. Kado uses a single source of truth (e.g., `cluster.yaml` and templates) for IaC external variables, which are passed to the IaC tools like Terraform and Ansible. 

This approach ensures that configurations are consistent across different environments and can be easily managed and updated.

## Overview

Kado uses a modular configuration approach where each module is called a "bead." Beads are blocks of configuration that define specific aspects of your infrastructure and can relay configurations to other beads. This document provides an overview of the current bead types, their functions, configured/allowed inputs, the processing flow, and the Kado file format including templates.

## Kado File Format

Kado uses `*.kd` files to define beads and their configurations. Each `.kd` file can contain multiple beads, and users can have as many `.kd` files and templates as needed. The structure of a `.kd` file is defined as follows:

### Basic Structure

A bead is defined using the following structure:

```hcl
#comment
bead "<name>" {
  key = "value"
  another_key = "another_value"
  #another_comment
}
```

- `<name>`: The name of the bead.
- `key = "value"`: Configuration settings.
- `#comment`: Comments for explanation.

### Example

```hcl
# This is an example bead configuration
bead "example_bead" {
  key1 = "value1"
  key2 = "value2"
  #additional comments
}
```

## Template Files

Kado uses templates to generate configuration files for various tools like Ansible and Terraform. The templates are stored in the `templates/` directory and can be customized as needed. Below are examples of custom template files used in Kado:

### Ansible Inventory Template

**Path**: `templates/ansible/inventory.tmpl`

```ini
<inventory.ini>
[proxmox]
{{join "proxmox.nodes.saathi01" "\n"}}
{{join "proxmox.nodes.saathi02" "\n"}}

[all:vars]
cluster_name={{.Get "proxmox.cluster_name"}}
ansible_user={{.Get "ansible.user"}}
ansible_python_interpreter={{.Get "ansible.python_interpreter"}}
```

### Terraform Variables Template

**Path**: `templates/terraform/vm.tfvars.tmpl`

```hcl
<vm.tfvars>
aws_region       = "{{.Get "aws.s3.region"}}"
pm_api_url       = "{{.Get "proxmox.api_url"}}"
pm_user          = "{{.Env "PM_USER"}}"
pm_password      = "{{.Env "PM_PASSWORD"}}"
vm_roles = {
  master       = {{.Get "proxmox.vm.roles.master"}}
  worker       = {{.Get "proxmox.vm.roles.worker"}}
  loadbalancer = {{.Get "proxmox.vm.roles.loadbalancer"}}
}
vm_template      = {{.Get "proxmox.vm.template"}}
vm_cpu           = {{.Get "proxmox.vm.cpu"}}
vm_memory        = {{.Get "proxmox.vm.memory"}}
vm_disk_size = "{{.Get "proxmox.vm.disk_size"}}"
vm_storage       = "{{.Get "proxmox.vm.storage"}}"
vm_network_bridge = "{{.Get "proxmox.vm.network_bridge"}}"
vm_network_model = "{{.Get "proxmox.vm.network_model"}}"
proxmox_nodes = {{ .GetKeysAsArray "proxmox.nodes" }}
ssh_public_key_content   = "/Users/janpreetsingh/.ssh/id_rsa.pub"
ssh_private_key          = "/Users/janpreetsingh/.ssh/id_rsa"
ssh_user  = "{{.Get "proxmox.vm.ssh_user"}}"
cloud_init_user_data_file = "templates/cloud_init_user_data.yaml"
k8s_master_setup_script  = "scripts/k8s_master_setup.sh"
k8s_worker_setup_script  = "scripts/k8s_worker_setup.sh"
haproxy_setup_script     = "scripts/haproxy_setup.sh"
haproxy_config_file      = "templates/haproxy.cfg"
s3_bucket                = "{{.Get "aws.s3.bucket"}}"
s3_key                   = "{{.Get "aws.s3.key"}}"
```

### Custom Template Functions

Kado provides custom template functions to enhance the templating capabilities:

- `Get`: Fetches the value of a specified key.
- `Env`: Fetches the value of an environment variable.
- `GetKeysAsArray`: Fetches the keys of a map as an array.

**Note**: The title of the output file (e.g., `<vm.tfvars>`) is added to the top of the file.

## Bead Types

### Ansible Bead

**Purpose**: Defines configurations for running Ansible playbooks.

**Configured/Allowed Inputs**:
- `enabled`: (boolean) Whether the Ansible bead is enabled.
- `source`: (string) Git repository URL for the Ansible playbook.
- `playbook`: (string) Path to the Ansible playbook.
- `extra_vars_file`: (boolean) Whether to use an extra variables file.
- `relay`: (string) Name of the bead to relay configurations to.
- `relay_field`: (string) Comma-separated list of key-value pairs to relay.

**Example**:
```hcl
bead "ansible" {
  enabled = false
  source = "git@github.com:janpreet/proxmox_ansible.git"
  playbook = "cluster.yaml"
  extra_vars_file = false
  relay = opa
  relay_field = "source=git@github.com:janpreet/proxmox_ansible.git,path=ansible/policies/proxmox.rego,input=ansible/cluster.yaml,package=data.proxmox.main.allowed"
  #extra_vars = "a=b"
}
```

### Terraform Bead

**Purpose**: Defines configurations for running Terraform.

**Configured/Allowed Inputs**:
- `enabled`: (boolean) Whether the Terraform bead is enabled.
- `source`: (string) Git repository URL for the Terraform configurations.
- `relay`: (string) Name of the bead to relay configurations to.
- `relay_field`: (string) Comma-separated list of key-value pairs to relay.

**Example**:
```hcl
bead "terraform" {
  source = "git@github.com:janpreet/proxmox_terraform.git"
  enabled = true
  relay = opa
  relay_field = "source=git@github.com:janpreet/proxmox_terraform.git,path=terraform/policies/proxmox.rego,input=terraform/plan.json,package=data.terraform.allow"
}
```

### OPA Bead

**Purpose**: Defines configurations for running Open Policy Agent (OPA) validations.

**Configured/Allowed Inputs**:
- `enabled`: (boolean) Whether the OPA bead is enabled.
- `path`: (string) Path to the OPA policy file.
- `input`: (string) Path to the input data file for OPA.
- `package`: (string) OPA package to evaluate.

**Example**:
```hcl
bead "opa" {
  enabled = true
  path = "path/to/opa/policy.rego"
  input = "path/to/opa/input.json"
  package = "data.example.allow"
}
```

### Custom Beads

**Purpose**: Define user-specific configurations.

**Configured/Allowed Inputs**: User-defined key-value pairs.

**Example**:
```hcl
bead "banana" {
  author = "Jane Doe"
  description = "This is a test bead"
  version = "3.1"
  status = "active"
}
```

## Processing Beads

### Basic Processing Flow

1. **Initialization**: Beads are read from `*.kd` files and stored in a list.
2. **Validation**: Each bead is validated based on its defined structure and required fields.
3. **Processing**: Beads are processed in the order they appear, executing their configurations.

### Relay Mechanism

A bead can relay its configurations to another bead using the `relay` and `relay_field` attributes. This mechanism allows one bead to pass its configurations to another bead for further processing.

**Relay Example**:
```hcl
bead "ansible" {
  enabled = false
  source = "git@github.com:janpreet/proxmox_ansible.git"
  playbook = "cluster.yaml"
  extra_vars_file = false
  relay = opa
  relay_field = "source=git@github.com:janpreet/proxmox_ansible.git,path=ansible/policies/proxmox.rego,input=ansible/cluster.yaml,package=data.proxmox.main.allowed"
}
```

### Relay Overrides

When a bead relays to another bead, it can override specific configurations using the `relay_field` attribute. The relay field is a comma-separated list of key-value pairs that specify the overrides.

### Processing Order

Beads are processed in the order they appear in the `.kd` files. If a bead relays to another bead, the relayed bead is processed next. This ensures that configurations are applied in a logical sequence.

### Preventing Duplicate Processing

A bead is processed once unless it is a relayed bead. To prevent duplicate processing:
- Keep track of processed beads using a map.
- Increment the count each time a bead is processed.
- Skip processing if the bead has already been processed, except for relayed beads.

**Example of Avoiding Duplicate Processing**:
```go
processed := make(map[string]int)

for _, b := range validBeads {
  if err := processBead(b, yamlData, beadMap, processed, &processedBeads, applyPlan, "", false); err != nil {
    log.Fatalf("Failed to process bead %s: %v", b.Name, err)


  }
}
```
### Cluster Configuration and Template Integration in Kado

Kado leverages a single source of truth file, typically named `cluster.yaml` or something relevant for ease of human readability, to drive the automation of Infrastructure as Code (IaC) using various beads. This configuration file is used to define all the necessary parameters and settings required for provisioning and managing infrastructure. Kado reads these configurations and uses them to populate templates that are then processed by different tools like Ansible, Terraform, and Terragrunt.

## Structure of example `cluster.yaml`

The `cluster.yaml` file follows a hierarchical structure, where different sections define specific configurations for various aspects of the infrastructure. Here's an example structure:

```yaml
kado:
  templates:
    - templates/ansible/inventory.tmpl
    - templates/terraform/backend.tfvars.tmpl
    - templates/terraform/vm.tfvars.tmpl

ansible:
  user: "user"
  python_interpreter: "/usr/bin/python3"

proxmox:
  cluster_name: "pmc"
  api_url: "https://1.2.3.4:8006/api2/json"
  user: "user"
  password: "password"
  nodes:
    saathi01:
      - 1.2.3.4
    saathi02:
      - 1.2.3.5
  vm:
    roles:
      master: 2
      worker: 3
      loadbalancer: 1
    template: 100
    cpu: 2
    memory: 2048
    storage: "local-lvm"
    disk_size: "10G"
    network_bridge: "vmbr0"
    network_model: "virtio"
    ssh_public_key_content: ""
    ssh_private_key: ""
    ssh_user: "ubuntu"

aws:
  s3:
    region: "aws-region"
    bucket: "s3-bucket"
    key: "tf-key"
```

### Key Sections

- **kado**: Defines the templates to be used for generating configuration files. Each template path is relative to the root of the project. This is the only section of yaml that needs to stay as is. Everything else is replacable key-value pairs.

## Using Templates in Kado

Kado processes the templates specified in the `kado.templates` section of `cluster.yaml` to generate the necessary configuration files. These templates use Go template syntax to dynamically populate values based on the `cluster.yaml` configurations.

### Example Templates

#### Ansible Inventory Template

**Path**: `templates/ansible/inventory.tmpl`

```hcl
<inventory.ini>
[proxmox]
{{join "proxmox.nodes.saathi01" "\n"}}
{{join "proxmox.nodes.saathi02" "\n"}}

[all:vars]
cluster_name={{.Get "proxmox.cluster_name"}}
ansible_user={{.Get "ansible.user"}}
ansible_python_interpreter={{.Get "ansible.python_interpreter"}}
```

#### Terraform Variables Template

**Path**: `templates/terraform/vm.tfvars.tmpl`

```hcl
<vm.tfvars>
aws_region       = "{{.Get "aws.s3.region"}}"
pm_api_url       = "{{.Get "proxmox.api_url"}}"
pm_user          = "{{.Env "PM_USER"}}"
pm_password      = "{{.Env "PM_PASSWORD"}}"
vm_roles = {
  master       = {{.Get "proxmox.vm.roles.master"}}
  worker       = {{.Get "proxmox.vm.roles.worker"}}
  loadbalancer = {{.Get "proxmox.vm.roles.loadbalancer"}}
}
vm_template      = {{.Get "proxmox.vm.template"}}
vm_cpu           = {{.Get "proxmox.vm.cpu"}}
vm_memory        = {{.Get "proxmox.vm.memory"}}
vm_disk_size = "{{.Get "proxmox.vm.disk_size"}}"
vm_storage       = "{{.Get "proxmox.vm.storage"}}"
vm_network_bridge = "{{.Get "proxmox.vm.network_bridge"}}"
vm_network_model = "{{.Get "proxmox.vm.network_model"}}"
proxmox_nodes = {{ .GetKeysAsArray "proxmox.nodes" }}
ssh_public_key_content   = "/path/to/id_rsa.pub"
ssh_private_key          = "/path/to/id_rsa"
ssh_user  = "{{.Get "proxmox.vm.ssh_user"}}"
cloud_init_user_data_file = "templates/cloud_init_user_data.yaml"
k8s_master_setup_script  = "scripts/k8s_master_setup.sh"
k8s_worker_setup_script  = "scripts/k8s_worker_setup.sh"
haproxy_setup_script     = "scripts/haproxy_setup.sh"
haproxy_config_file      = "templates/haproxy.cfg"
s3_bucket                = "{{.Get "aws.s3.bucket"}}"
s3_key                   = "{{.Get "aws.s3.key"}}"
```

## Driving Bead Automation with `cluster.yaml`

The `cluster.yaml` file serves as the single source of truth for all configurations, driving the automation process within Kado. Each bead in Kado processes its respective templates and configuration settings as defined in `cluster.yaml`.

### Processing Flow

1. **Read `cluster.yaml`**: Kado reads the `cluster.yaml` file to gather all configurations.
2. **Load Templates**: The templates defined in the `kado.templates` section are loaded.
3. **Process Beads**: Each bead processes its templates and executes the necessary commands. The templates are populated with values from `cluster.yaml`.
4. **Relay to OPA**: If a bead is configured to relay to OPA, the generated plan (e.g., Terraform or Terragrunt plan) is evaluated by OPA before proceeding with the apply step.

By defining configurations in `cluster.yaml` and using Kado's templating system, users can achieve a seamless and automated workflow for managing their infrastructure.