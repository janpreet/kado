// File: packages/helper/processor.go

package helper

import (
	"fmt"
	"path/filepath"

	"github.com/janpreet/kado/packages/bead"
	"github.com/janpreet/kado/packages/config"
	"github.com/janpreet/kado/packages/engine"
	"github.com/janpreet/kado/packages/opa"
	"github.com/janpreet/kado/packages/render"
	"github.com/janpreet/kado/packages/terraform"
	"github.com/janpreet/kado/packages/terragrunt"
)

func ProcessAnsibleBead(b bead.Bead, yamlData map[string]interface{}, relayToOPA bool, applyPlan bool) error {
	fmt.Println("Processing Ansible templates...")
	templatePaths, ok := yamlData["kado"].(map[string]interface{})["templates"].([]interface{})
	if !ok {
		return fmt.Errorf("no templates defined for Ansible in the YAML configuration")
	}
	err := render.ProcessTemplates(convertTemplatePaths(templatePaths), yamlData)
	if err != nil {
		return fmt.Errorf("failed to process Ansible templates: %v", err)
	}
	if relayToOPA {
		fmt.Println("Ansible bead is relayed to OPA for evaluation.")
	}
	if playbook, ok := b.Fields["playbook"]; ok && playbook != "" {
		playbookPath := filepath.Join(config.LandingZone, b.Name, playbook)
		inventoryPath := b.Fields["inventory"]
		if inventoryPath == "" {
			inventoryPath = filepath.Join(config.LandingZone, "inventory.ini")
		}
		extraVarsFile := false
		if extraVarsFileFlag, ok := b.Fields["extra_vars_file"]; ok && extraVarsFileFlag == "true" {
			extraVarsFile = true
		}
		fmt.Printf("Running Ansible playbook: %s with inventory: %s\n", playbookPath, inventoryPath)
		if !FileExists(playbookPath) {
			return fmt.Errorf("playbook file does not exist: %s", playbookPath)
		}
		if !relayToOPA || (relayToOPA && applyPlan) {
			err := engine.HandleAnsible(b, convertYAMLToSlice(yamlData), extraVarsFile)
			if err != nil {
				return fmt.Errorf("failed to run Ansible: %v", err)
			}
		} else {
			fmt.Println("Skipping Ansible playbook apply due to OPA evaluation or missing 'set' flag.")
		}
	}
	return nil
}

func ProcessTerraformBead(b bead.Bead, yamlData map[string]interface{}, applyPlan bool) error {
	fmt.Println("Processing Terraform templates...")
	templatePaths, ok := yamlData["kado"].(map[string]interface{})["templates"].([]interface{})
	if !ok {
		return fmt.Errorf("no templates defined for Terraform in the YAML configuration")
	}
	err := render.ProcessTemplates(convertTemplatePaths(templatePaths), yamlData)
	if err != nil {
		return fmt.Errorf("failed to process Terraform templates: %v", err)
	}
	fmt.Println("Running Terraform plan...")
	err = terraform.HandleTerraform(b, config.LandingZone, applyPlan)
	if err != nil {
		return fmt.Errorf("failed to run Terraform: %v", err)
	}
	return nil
}

func ProcessOPABead(b bead.Bead, applyPlan bool, originBead string) error {
	fmt.Println("Processing OPA validation...")
	fmt.Printf("DEBUG: Calling HandleOPA with originBead: %s\n", originBead)
	err := opa.HandleOPA(b, config.LandingZone, applyPlan, originBead)
	if err != nil {
		return fmt.Errorf("failed to process OPA: %v", err)
	}
	return nil
}

func ProcessTerragruntBead(b bead.Bead, yamlData map[string]interface{}, applyPlan bool) error {
	fmt.Println("Processing Terragrunt templates...")
	templatePaths, ok := yamlData["kado"].(map[string]interface{})["templates"].([]interface{})
	if !ok {
		return fmt.Errorf("no templates defined for Terragrunt in the YAML configuration")
	}
	err := render.ProcessTemplates(convertTemplatePaths(templatePaths), yamlData)
	if err != nil {
		return fmt.Errorf("failed to process Terragrunt templates: %v", err)
	}
	fmt.Println("Running Terragrunt plan...")
	err = terragrunt.HandleTerragrunt(b, config.LandingZone, applyPlan)
	if err != nil {
		return fmt.Errorf("failed to run Terragrunt: %v", err)
	}
	return nil
}

// Helper functions

func convertTemplatePaths(paths []interface{}) []string {
	var result []string
	for _, path := range paths {
		if strPath, ok := path.(string); ok {
			result = append(result, strPath)
		}
	}
	return result
}

func convertYAMLToSlice(yamlData map[string]interface{}) []map[string]interface{} {
	return []map[string]interface{}{yamlData}
}