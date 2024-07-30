[![Build, Docker, and Release](https://github.com/janpreet/kado/actions/workflows/docker-release.yaml/badge.svg)](https://github.com/janpreet/kado/actions/workflows/docker-release.yaml)[![Sensitive Data Check](https://github.com/janpreet/kado/actions/workflows/sensitive-data-check.yml/badge.svg)](https://github.com/janpreet/kado/actions/workflows/sensitive-data-check.yml)
<p align="center"><img src="https://raw.githubusercontent.com/janpreet/kado/main/assets/kado_dark.png" data-canonical-src="https://raw.githubusercontent.com/janpreet/kado/main/assets/kado_dark.png" width="300" height="250" /></p>

## Introduction

Kado is a modular configuration management tool designed to streamline and automate the provisioning and configuration of infrastructure using tools like Ansible, Terraform, and Terragrunt. It provides a flexible framework for defining and processing configurations through a concept called "beads," which are modular units of configuration.

## Table of Contents

- [Overview](#overview)
- [Configuration Files](#configuration-files)
  - [cluster.yaml](#clusteryaml)
  - [Template Files](#template-files)
- [Beads](#beads)
  - [Ansible Bead](#ansible-bead)
  - [Terraform Bead](#terraform-bead)
  - [OPA Bead](#opa-bead)
  - [Terragrunt Bead](#terragrunt-bead)
- [Usage](#usage)
  - [Commands](#commands)
  - [Getting Started](#getting-started)
  - [Configuration](#configuration)
- [Upcoming Improvements](#upcoming-improvements)
- [Code of Conduct](#code-of-conduct)

## Overview

Kado is a Bring Your Own Code (BYOC) tool that leverages your existing Ansible, Terraform, and Terragrunt configurations, and provides a single source of truth for your infrastructure parameters. It uses `*.kd` files for defining beads and `*.yaml` for centralized configuration, making it easy to manage and relay configurations across different infrastructure components.

## Configuration Files

### cluster.yaml

The `cluster.yaml` file serves as the single source of truth for your infrastructure configuration. It contains various parameters that are used to drive the automation of infrastructure provisioning and configuration.

Example `cluster.yaml`:

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

### Template Files

Template files are used to generate configuration files for various tools like Ansible and Terraform. These templates are stored in the `templates/` directory and can be customized as needed to meet the specific needs of your infrastructure, providing flexibility and control over the configuration process.

Example `vm.tfvars.tmpl`:

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
vm_storage       = "{{.Get "proxmox.vm.storage"}}
vm_network_bridge = "{{.Get "proxmox.vm.network_bridge"}}
vm_network_model = "{{.Get "proxmox.vm.network_model"}}
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

## Beads

Beads are modular units of configuration in Kado. Each bead defines specific aspects of your infrastructure and can relay configurations to other beads. Kado uses `*.kd` files to define beads and their configurations. Users can have as many `.kd` files and templates as needed, allowing for a highly customizable and scalable setup.

### Ansible Bead

**Purpose**: Defines configurations for running Ansible playbooks.

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

**Example**:

```hcl
bead "opa" {
  enabled = true
  path = "path/to/opa/policy.rego"
  input = "path/to/opa/input.json"
  package = "data.example.allow"
}
```

### Terragrunt Bead

**Purpose**: Defines configurations for running Terragrunt.

**Example**:

```hcl
bead "terragrunt" {
  source = "git@github.com:janpreet/proxmox_terragrunt.git"
  enabled = true
  relay = opa
  relay_field = "source=git@github.com:janpreet/proxmox_terragrunt.git,path=terragrunt/policies/proxmox.rego,input=terragrunt/plan.json,package=data.terraform.allow"
}
```

## Usage

### Commands

- `kado [file.yaml]`: Runs the default configuration and processing of beads. You may pass a specific YAML file to Kado. If no file is specified, Kado scans all YAML files in the current directory.
- `kado set`: Applies the configuration and processes beads with the `set` flag.
- `kado fmt [dir]`: Formats `.kd` files in the specified directory.
- `kado ai`: Runs AI-based recommendations if enabled in the `~/.kdconfig` configuration.
- `kado config`: Displays the current configuration and order of execution.

### Getting Started

1. **Download the latest release** from GitHub.
2. **Create your configuration files** (`cluster.yaml` and `.kd` files).
3. **Define your templates** in the `templates/` directory.
4. **Run Kado** using one of the commands listed above.

### Configuration

Create a `.kdconfig` file in your home directory to enable AI recommendations:

```plaintext
AI_API_KEY=<your_api_key>
AI_MODEL=gpt-3.5-turbo
AI_CLIENT=chatgpt
AI_ENABLED=true
```

### Outputs

- **Processed Beads**: Lists the beads that have been successfully processed.
- **Skipped Beads**: Lists the beads that were skipped and the reasons for skipping.

## Conclusion

Kado aims to simplify and streamline the management of your infrastructure as code. By providing a modular, consistent, and automated framework, Kado helps you reduce complexity, minimize errors, and achieve efficient infrastructure management. Whether you are provisioning resources with Terraform, managing configurations with Ansible, or enforcing policies with OPA, Kado brings everything together into a cohesive and powerful tool.

## Upcoming Improvements

- More tests and better test coverage.
- Support for CDK and Pulumi.
- Improved error handling and logging.
- More customizable and dynamic templating functions.

Dive into the Kado project and experience a new level of simplicity and efficiency in managing your infrastructure!
[Configuration](https://github.com/janpreet/kado/blob/main/assets/Configuration.md), [How to](https://github.com/janpreet/kado/blob/main/assets/How-to.md), [Structure](https://github.com/janpreet/kado/blob/main/assets/Structure.md), [TLDR](https://github.com/janpreet/kado/blob/main/assets/TLDR.md)