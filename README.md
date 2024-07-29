[![Build, Docker, and Release](https://github.com/janpreet/kado/actions/workflows/docker-release.yaml/badge.svg)](https://github.com/janpreet/kado/actions/workflows/docker-release.yaml)[![Sensitive Data Check](https://github.com/janpreet/kado/actions/workflows/sensitive-data-check.yml/badge.svg)](https://github.com/janpreet/kado/actions/workflows/sensitive-data-check.yml)
<p align="center"><img src="https://raw.githubusercontent.com/janpreet/kado/main/assets/kado_dark.png" data-canonical-src="https://raw.githubusercontent.com/janpreet/kado/main/assets/kado_dark.png" width="300" height="250" /></p>

## Introduction

Kado is a powerful and flexible tool designed to streamline the management of infrastructure configurations using beads through a modular and declarative approach. Whether you're managing virtual machines, Kubernetes clusters, or other infrastructure components, Kado provides a cohesive framework to integrate and automate your Infrastructure as Code (IaC) processes using Terraform, Ansible, and Open Policy Agent (OPA).

## Problems Kado Solves

### 1. **Consistency in Infrastructure Configurations**

Maintaining consistency in infrastructure configurations across different environments can be challenging. Inconsistent configurations can lead to unexpected behaviors, security vulnerabilities, and operational inefficiencies.

**How Kado Solves It:**
- Kado uses a single source of truth for infrastructure configurations, defined in `*.kd` files and templates. This ensures that configurations are consistent and repeatable across all environments.

### 2. **Modular and Scalable Configuration Management**

As infrastructure grows in complexity, managing configurations becomes increasingly difficult. Large monolithic configuration files are hard to maintain and scale.

**How Kado Solves It:**
- Kado introduces a modular approach with "beads," which are blocks of configurations that define specific aspects of your infrastructure. This modularity allows you to manage, update, and scale your configurations easily.

### 3. **Automation and Integration of IaC Tools**

Integrating and automating different IaC tools like Terraform and Ansible can be cumbersome. Each tool has its own syntax, workflows, and integration points, which can complicate automation.

**How Kado Solves It:**
- Kado seamlessly integrates Terraform and Ansible by using templates and a unified configuration structure. It automates the execution of these tools, ensuring smooth and efficient workflows.

### 4. **Policy Enforcement and Compliance**

Ensuring that your infrastructure complies with security policies and operational guidelines is critical. However, manually enforcing these policies can be error-prone and time-consuming.

**How Kado Solves It:**
- Kado incorporates Open Policy Agent (OPA) to enforce policies and compliance checks. By relaying configurations to OPA, Kado ensures that only approved policies are applied, enhancing security and compliance.

### 5. **Simplified Configuration Management**

Managing infrastructure configurations often requires deep knowledge of various tools and their configurations. This complexity can slow down development and operations teams.

**How Kado Solves It:**
- Kado abstracts the complexities of individual IaC tools and provides a simplified, unified interface for managing configurations. This reduces the learning curve and enables teams to focus on their core tasks.

## How Kado Achieves This

### 1. **Declarative Configuration Files**

Kado uses `*.kd` files to define infrastructure configurations in a declarative manner. These files contain "beads," which are modular blocks that specify configurations for different aspects of your infrastructure.

### 2. **Templating System**

Kado leverages a powerful templating system to generate configuration files for Terraform and Ansible. Templates can be customized to fit specific needs, ensuring flexibility and adaptability.

### 3. **Automation of IaC Workflows**

Kado automates the execution of Terraform and Ansible by processing the `*.kd` files and applying the configurations. This automation streamlines workflows and reduces the potential for human error.

### 4. **Policy Enforcement with OPA**

By integrating Open Policy Agent, Kado enforces policies and compliance checks on your infrastructure configurations. Beads can relay their configurations to OPA, ensuring that only compliant configurations are applied.

### 5. **Modular and Extensible Design**

Kado's bead structure allows for modular and extensible configuration management. Users can define custom beads and templates to fit their unique infrastructure requirements.

### 6. **Bring Your Own Code (BYOC)**

Kado supports BYOC, allowing users to plug in their own parameterized IaC code. This means you can bring your existing Terraform and Ansible configurations and integrate them into Kado.

### 7. Centralized Variable Management
Kado uses a central configuration file (e.g., cluster.yaml) to store variable values that are used across different beads and templates. This ensures consistency and simplifies the management of configuration variables.


## Key Features

### Modular Configuration

Kado uses a bead-based configuration system where each bead represents a distinct aspect of your infrastructure. This modular approach allows you to define specific configurations for different tools and components, making it easy to manage and update your infrastructure as needed.

### Single Source of Truth

Kado leverages a single source of truth for configuration variables, ensuring consistency across different environments. By defining variables in a centralized configuration file (e.g., `cluster.yaml`), Kado ensures that all infrastructure components use the same set of parameters, reducing the risk of configuration drift.

### Integration with Popular Tools

Kado seamlessly integrates with popular infrastructure management tools like Terraform, Ansible, and OPA. This integration allows you to harness the power of these tools while benefiting from Kado's unified configuration and management framework.

### Automation and Relay Mechanism

Kado automates the deployment and management of infrastructure by processing beads in a defined order. Beads can relay their configurations to other beads, enabling a flexible and dynamic workflow. This relay mechanism ensures that configurations are applied logically and consistently across different components.

### Policy Enforcement with OPA

Kado supports policy enforcement using OPA. If OPA is enabled and configured, Kado ensures that infrastructure changes comply with defined policies before applying them. This feature helps maintain security and compliance standards across your infrastructure.

## Getting Started

### Easy Configuration

Kado uses `*.kd` files to define beads and their configurations. Users can have as many `.kd` files and templates as needed, allowing for a highly customizable and scalable setup.

### Flexible Templates

Kado supports custom templates for generating configuration files for various tools. These templates can be tailored to meet the specific needs of your infrastructure, providing flexibility and control over the configuration process.

### Simplified Deployment

Kado simplifies the deployment process by automating the execution of configured beads. Whether you are using Terraform for provisioning resources, Ansible for configuration management, or OPA for policy enforcement, Kado ensures that everything works together seamlessly.

### Example Workflow

1. **Define Beads**: Create `*.kd` files to define your infrastructure components using beads.
2. **Configure Templates**: Customize templates to generate the necessary configuration files for your tools.
3. **Run Kado**: Use Kado commands to process and apply your configurations, ensuring that your infrastructure is deployed and managed consistently.

## Conclusion

Kado aims to simplify and streamline the management of your infrastructure as code. By providing a modular, consistent, and automated framework, Kado helps you reduce complexity, minimize errors, and achieve efficient infrastructure management. Whether you are provisioning resources with Terraform, managing configurations with Ansible, or enforcing policies with OPA, Kado brings everything together into a cohesive and powerful tool.

Dive into the Kado project and experience a new level of simplicity and efficiency in managing your infrastructure!
[Configuration](https://github.com/janpreet/kado/blob/main/assets/Configuration.md), [How to](https://github.com/janpreet/kado/blob/main/assets/How-to.md), [Structure](https://github.com/janpreet/kado/blob/main/assets/Structure.md)
