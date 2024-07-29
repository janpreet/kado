package display

import (
	"fmt"
	"github.com/janpreet/kado/packages/bead"
)

func DisplayBeads(kdBeads map[string]bead.Bead, parsedYAMLs []map[string]interface{}) {
	for _, b := range kdBeads {
		fmt.Printf("Bead: %s\n", b.Name)
		for k, v := range b.Fields {
			fmt.Printf("  %s = %s\n", k, v)
		}
		if b.Name == "ansible" {
			fmt.Printf("Bead details - Name: %s, Playbook: %s, Inventory: %s, Source: %s, ExtraVarsFile: %s\n", b.Name, b.Fields["playbook"], b.Fields["inventory"], b.Fields["source"], "LandingZone/extra_vars.yaml")
		}
	}
}

func DisplayYAMLs(parsedYAMLs []map[string]interface{}) {
	for _, yamlContent := range parsedYAMLs {
		DisplayYAML(yamlContent)
	}
}

func DisplayYAML(yamlContent map[string]interface{}) {
	fmt.Println("YAML Content:")
	for key, value := range yamlContent {
		fmt.Printf("  %s = %v\n", key, value)
	}
}

func DisplayTemplateOutput(outputPath string) {
	fmt.Printf("Template processed successfully. Output written to: %s\n", outputPath)
}

func DisplayBead(b bead.Bead) {
	fmt.Printf("Bead: %s\n", b.Name)
	for key, value := range b.Fields {
		fmt.Printf("  %s = %s\n", key, value)
	}
}

func DisplayBeadConfig(beads []bead.Bead) {
	fmt.Println("Bead Configuration and Order of Execution:")

	displayed := make(map[string]bool)

	var displayBeadChain func(string)
	displayBeadChain = func(name string) {
		if displayed[name] {
			return
		}
		for _, b := range beads {
			if b.Name == name {
				if displayed[name] {
					continue
				}
				if len(displayed) > 0 {
					fmt.Println("↓")
				}
				fmt.Printf("Bead: %s\n", b.Name)
				for key, value := range b.Fields {
					fmt.Printf("  %s = %s\n", key, value)
				}
				displayed[name] = true
				if relay, ok := b.Fields["relay"]; ok {
					displayBeadChain(relay)
				}
				break
			}
		}
	}

	for _, b := range beads {
		if !displayed[b.Name] {
			if len(displayed) > 0 {
				fmt.Println()
			}
			displayBeadChain(b.Name)
		}
	}
}
