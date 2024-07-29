package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/janpreet/kado/packages/bead"
	"github.com/janpreet/kado/packages/engine"
	"github.com/janpreet/kado/packages/terraform"
	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v3"
)

func HandleOPA(b bead.Bead, landingZone string, applyPlan bool, originBead string) error {
	fmt.Printf("Processing OPA bead:\n")
	for key, val := range b.Fields {
		fmt.Printf("  %s = %s\n", key, val)
	}

	inputPath, ok := b.Fields["input"]
	if !ok {
		return fmt.Errorf("input path not specified in bead")
	}
	fullInputPath := filepath.Join(landingZone, inputPath)
	fmt.Printf("Reading input file from path: %s\n", fullInputPath)
	inputData, err := os.ReadFile(fullInputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	var input interface{}
	if filepath.Ext(fullInputPath) == ".yaml" || filepath.Ext(fullInputPath) == ".yml" {
		if err := yaml.Unmarshal(inputData, &input); err != nil {
			return fmt.Errorf("failed to unmarshal YAML input file: %v", err)
		}
	} else {
		if err := json.Unmarshal(inputData, &input); err != nil {
			return fmt.Errorf("failed to unmarshal JSON input file: %v", err)
		}
	}

	policyPath, ok := b.Fields["path"]
	if !ok {
		return fmt.Errorf("policy path not specified in bead")
	}
	fullPolicyPath := filepath.Join(landingZone, policyPath)
	fmt.Printf("Reading policy file from path: %s\n", fullPolicyPath)
	policyData, err := os.ReadFile(fullPolicyPath)
	if err != nil {
		return fmt.Errorf("failed to read policy file: %v", err)
	}

	packageQuery := "data.terraform.allow"
	if pkg, ok := b.Fields["package"]; ok {
		packageQuery = pkg
	}
	fmt.Printf("Evaluating package: %s\n", packageQuery)

	ctx := context.Background()
	query, err := rego.New(
		rego.Query(packageQuery),
		rego.Module("policy.rego", string(policyData)),
	).PrepareForEval(ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare rego query: %v", err)
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("failed to evaluate rego query: %v", err)
	}

	if len(results) == 0 || len(results[0].Expressions) == 0 || results[0].Expressions[0].Value != true {
		fmt.Println("Input is denied by OPA policy.")
		if applyPlan {
			fmt.Println("Skipping action because the input was denied.")
		}
	} else {
		fmt.Println("Input is allowed by OPA policy.")
		if applyPlan {
			switch originBead {
			case "terraform":
				fmt.Println("Applying terraform plan...")
				err = terraform.HandleTerraform(b, landingZone, true)
				if err != nil {
					return fmt.Errorf("failed to apply terraform plan: %v", err)
				}
			case "ansible":
				fmt.Println("Applying ansible playbook...")
				err = handleAnsibleRelay(b, landingZone)
				if err != nil {
					return fmt.Errorf("failed to run Ansible: %v", err)
				}
			default:
				fmt.Println("Skipping apply action because origin bead is not terraform or ansible.")
			}
		} else {
			fmt.Println("Skipping apply action because 'set' was not passed.")
		}
	}

	return nil
}

func convertYAMLToSlice(yamlData map[string]interface{}) []map[string]interface{} {
	result := []map[string]interface{}{yamlData}
	return result
}

func handleAnsibleRelay(b bead.Bead, landingZone string) error {

	yamlPath := filepath.Join(landingZone, "ansible", "cluster.yaml")
	yamlData, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to read YAML config: %v", err)
	}

	var yamlContent map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &yamlContent); err != nil {
		return fmt.Errorf("failed to unmarshal YAML config: %v", err)
	}

	extraVarsFile := false
	if extraVarsFileFlag, ok := b.Fields["extra_vars_file"]; ok && extraVarsFileFlag == "true" {
		extraVarsFile = true
	}

	err = engine.HandleAnsible(b, convertYAMLToSlice(yamlContent), extraVarsFile)
	if err != nil {
		return fmt.Errorf("failed to run Ansible: %v", err)
	}

	return nil
}
