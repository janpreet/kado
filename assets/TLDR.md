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

Kado automates the execution of Terraform/ Terragrunt and Ansible by processing the `*.kd` files and applying the configurations. This automation streamlines workflows and reduces the potential for human error.

### 4. **Policy Enforcement with OPA**

By integrating Open Policy Agent, Kado enforces policies and compliance checks on your infrastructure configurations. Beads can relay their configurations to OPA, ensuring that only compliant configurations are applied.

### 5. **Modular and Extensible Design**

Kado's bead structure allows for modular and extensible configuration management. Users can define custom beads and templates to fit their unique infrastructure requirements.

### 6. **Bring Your Own Code (BYOC)**

Kado supports BYOC, allowing users to plug in their own parameterized IaC code. This means you can bring your existing Terraform and Ansible configurations and integrate them into Kado.

### 7. Centralized Variable Management
Kado uses a central configuration file (e.g., cluster.yaml) to store variable values that are used across different beads and templates. This ensures consistency and simplifies the management of configuration variables.