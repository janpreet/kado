package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/janpreet/kado/packages/bead"
	"github.com/janpreet/kado/packages/config"
	"github.com/janpreet/kado/packages/display"
	"github.com/janpreet/kado/packages/engine"
	"github.com/janpreet/kado/packages/helper"
	"github.com/janpreet/kado/packages/opa"
	"github.com/janpreet/kado/packages/render"
	"github.com/janpreet/kado/packages/terraform"
	"github.com/janpreet/kado/packages/terragrunt" // Import the new package
)

func convertYAMLToSlice(yamlData map[string]interface{}) []map[string]interface{} {
	result := []map[string]interface{}{yamlData}
	return result
}

func applyRelayOverrides(b *bead.Bead) map[string]string {
	overrides := make(map[string]string)
	if relayField, ok := b.Fields["relay_field"]; ok {
		pairs := strings.Split(relayField, ",")
		for _, pair := range pairs {
			keyValue := strings.SplitN(pair, "=", 2)
			if len(keyValue) == 2 {
				overrides[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
			}
		}
	}
	return overrides
}

func processBead(b bead.Bead, yamlData map[string]interface{}, beadMap map[string]bead.Bead, processed map[string]int, processedBeads *[]string, applyPlan bool, originBead string, relayToOPA bool) error {
	if count, ok := processed[b.Name]; ok && count > 0 && originBead == "" {
		return nil
	}

	fmt.Printf("Processing bead: %s\n", b.Name)

	if originBead != "" {
		repoPath := filepath.Join(config.LandingZone, b.Name)
		if helper.FileExists(repoPath) {
			fmt.Printf("Removing existing repository at: %s\n", repoPath)
			err := os.RemoveAll(repoPath)
			if err != nil {
				return fmt.Errorf("failed to remove existing repository for bead %s: %v", b.Name, err)
			}
		}
	}

	if source, ok := b.Fields["source"]; ok && source != "" {
		refs := ""
		if refsVal, ok := b.Fields["refs"]; ok {
			refs = refsVal
		}
		err := helper.CloneRepo(source, config.LandingZone, b.Name, refs)
		if err != nil {
			return fmt.Errorf("failed to clone repo for bead %s: %v", b.Name, err)
		}
	}

	display.DisplayBead(b)

	if b.Name == "ansible" {
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
			if !helper.FileExists(playbookPath) {
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
	}

	if b.Name == "terraform" {
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
	}

	if b.Name == "opa" {
		fmt.Println("Processing OPA validation...")
		err := opa.HandleOPA(b, config.LandingZone, applyPlan, originBead)
		if err != nil {
			return fmt.Errorf("failed to process OPA: %v", err)
		}
	}

	if b.Name == "terragrun" {
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
	}

	*processedBeads = append(*processedBeads, b.Name)
	processed[b.Name]++

	if relay, ok := b.Fields["relay"]; ok {
		if relayBead, ok := beadMap[relay]; ok {

			overrides := applyRelayOverrides(&b)
			for key, value := range overrides {
				relayBead.Fields[key] = value
			}
			return processBead(relayBead, yamlData, beadMap, processed, processedBeads, applyPlan, b.Name, b.Name == "opa")
		}
	}

	return nil
}

func convertTemplatePaths(paths []interface{}) []string {
	var result []string
	for _, path := range paths {
		if strPath, ok := path.(string); ok {
			result = append(result, strPath)
		}
	}
	return result
}

func main() {
	var yamlFilePath string
	if len(os.Args) > 1 && strings.HasSuffix(os.Args[1], ".yaml") {
		yamlFilePath = os.Args[1]
	} else {
		yamlFilePath = "cluster.yaml"
	}

	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("Version:", config.Version)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "config" {
		kdFiles, err := render.GetKDFiles(".")
		if err != nil {
			log.Fatalf("Failed to get KD files: %v", err)
		}

		var beads []bead.Bead
		for _, kdFile := range kdFiles {
			bs, err := config.LoadBeadsConfig(kdFile)
			if err != nil {
				log.Fatalf("Failed to load beads config from %s: %v", kdFile, err)
			}
			beads = append(beads, bs...)
		}

		display.DisplayBeadConfig(beads)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "fmt" {
		dir := "."
		if len(os.Args) > 2 {
			dir = os.Args[2]
		}
		err := engine.FormatKDFilesInDir(dir)
		if err != nil {
			log.Fatalf("Error formatting .kd files: %v", err)
		}
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "ai" {
		engine.RunAI()
		return
	}

	fmt.Println("Starting processing-")

	applyPlan := len(os.Args) > 1 && os.Args[1] == "set"

	kdFiles, err := render.GetKDFiles(".")
	if err != nil {
		log.Fatalf("Failed to get KD files: %v", err)
	}

	var beads []bead.Bead
	for _, kdFile := range kdFiles {
		bs, err := config.LoadBeadsConfig(kdFile)
		if err != nil {
			log.Fatalf("Failed to load beads config from %s: %v", kdFile, err)
		}
		beads = append(beads, bs...)
	}

	yamlData, err := config.LoadYAMLConfig(yamlFilePath)
	if err != nil {
		log.Fatalf("Failed to load YAML config: %v", err)
	}

	err = helper.SetupLandingZone()
	if err != nil {
		log.Fatalf("Failed to setup LandingZone: %v", err)
	}

	var invalidBeadNames []string
	var processedBeads []string

	validBeads, invalidBeadReasons := config.GetValidBeadsWithDefaultEnabled(beads)

	beadMap := make(map[string]bead.Bead)
	for _, b := range validBeads {
		beadMap[b.Name] = b
	}

	processed := make(map[string]int)

	for _, b := range validBeads {
		if err := processBead(b, yamlData, beadMap, processed, &processedBeads, applyPlan, "", false); err != nil {
			log.Fatalf("Failed to process bead %s: %v", b.Name, err)
		}
	}

	for beadIndex, reason := range invalidBeadReasons {
		beadName := fmt.Sprintf("bead_%d", beadIndex)
		fmt.Printf("Skipping bead: %s, Reason: %s\n", beadName, reason)
		invalidBeadNames = append(invalidBeadNames, fmt.Sprintf("%s: %s", beadName, reason))
	}

	if len(processedBeads) > 0 {
		fmt.Println("\nProcessed beads:")
		for _, name := range processedBeads {
			fmt.Printf("  - %s\n", name)
		}
	}

	if len(invalidBeadNames) > 0 {
		fmt.Println("\nSkipped beads:")
		for _, name := range invalidBeadNames {
			fmt.Printf("  - %s\n", name)
		}
	}
}
