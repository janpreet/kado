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
	"github.com/janpreet/kado/packages/render"
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
	config.DebugPrint("DEBUG: processBead called for %s (Enabled: %v, Origin: %s)\n", b.Name, *b.Enabled, originBead)
    
    if b.Enabled != nil && !*b.Enabled {
		config.DebugPrint("DEBUG: Skipping disabled bead: %s\n", b.Name)       
		return nil
    }

    if count, ok := processed[b.Name]; ok && count > 0 && originBead == "" {
		config.DebugPrint("DEBUG: Skipping already processed bead: %s\n", b.Name)
        return nil
    }

	config.DebugPrint("DEBUG: Actually processing bead: %s\n", b.Name)

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

    switch b.Name {
    case "ansible":
        err := helper.ProcessAnsibleBead(b, yamlData, relayToOPA, applyPlan)
        if err != nil {
            return err
        }
    case "terraform":
        err := helper.ProcessTerraformBead(b, yamlData, applyPlan)
        if err != nil {
            return err
        }
    case "opa":
        err := helper.ProcessOPABead(b, applyPlan, originBead)
        if err != nil {
            return err
        }
    case "terragrunt":
        err := helper.ProcessTerragruntBead(b, yamlData, applyPlan)
        if err != nil {
            return err
        }
    default:
        return fmt.Errorf("unknown bead type: %s", b.Name)
    }

    processed[b.Name]++
    *processedBeads = append(*processedBeads, b.Name)
    config.DebugPrint("DEBUG: Added %s to processedBeads\n", b.Name)

    if relay, ok := b.Fields["relay"]; ok {
		config.DebugPrint("DEBUG: Relay found for %s to %s\n", b.Name, relay)
        if relayBead, ok := beadMap[relay]; ok {
            if relayBead.Enabled != nil && !*relayBead.Enabled {
                config.DebugPrint("DEBUG: Skipping disabled relay bead: %s\n", relayBead.Name)
                return nil
            }
            overrides := applyRelayOverrides(&b)
            for key, value := range overrides {
                relayBead.Fields[key] = value
            }
            config.DebugPrint("DEBUG: Calling processBead for relay %s\n", relayBead.Name)
            return processBead(relayBead, yamlData, beadMap, processed, processedBeads, applyPlan, b.Name, b.Name == "opa")
        } else {
            config.DebugPrint("DEBUG: Relay bead %s not found in beadMap\n", relay)
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

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Println("Version:", config.Version)
			return

		case "config":
			handleConfigCommand()
			return

		case "fmt":
			handleFormatCommand()
			return

		case "ai":
			engine.RunAI()
			return

		case "keybase":
			if len(os.Args) < 3 {
				fmt.Println("Usage: kado keybase [debug] <command>")
				return
			}
			helper.HandleKeybaseCommand(os.Args[2:])
			return

		default:
			if strings.HasSuffix(os.Args[1], ".yaml") {
				yamlFilePath = os.Args[1]
			} else {
				yamlFilePath = "cluster.yaml"
			}
		}
	} else {
		yamlFilePath = "cluster.yaml"
	}

	fmt.Println("Starting processing-")

	applyPlan := len(os.Args) > 1 && os.Args[1] == "set"

	kdFiles, err := render.GetKDFiles(".")
	if err != nil {
		log.Fatalf("Failed to get KD files: %v", err)
	}

	beadMap := make(map[string]bead.Bead)
	var primaryKdFile string
	
	for i, kdFile := range kdFiles {
		config.DebugPrint("DEBUG: Loading file: %s\n", kdFile)
		bs, err := config.LoadBeadsConfig(kdFile)
		if err != nil {
			log.Fatalf("Failed to load beads config from %s: %v", kdFile, err)
		}
		
		if i == 0 {
			primaryKdFile = kdFile
		}
		
		for _, b := range bs {
			if _, ok := beadMap[b.Name]; ok {
				if kdFile != primaryKdFile {
					fmt.Printf("WARNING: Ignoring conflicting configuration for bead %s in file %s. Using configuration from %s\n", b.Name, kdFile, primaryKdFile)
				} else {
					beadMap[b.Name] = b
					config.DebugPrint("DEBUG: Updated bead %s (Enabled: %v) from primary file %s\n", b.Name, *b.Enabled, kdFile)
				}
			} else {
				beadMap[b.Name] = b
				config.DebugPrint("DEBUG: Loaded new bead %s (Enabled: %v) from file %s\n", b.Name, *b.Enabled, kdFile)
			}
		}
	}
	
	var allBeads []bead.Bead
	for _, b := range beadMap {
		allBeads = append(allBeads, b)
	}
	
	validBeads, invalidBeadReasons := config.GetValidBeadsWithDefaultEnabled(allBeads)
	
	config.DebugPrint("DEBUG: Final bead configurations:")
	for _, b := range validBeads {
		fmt.Printf("  - %s (Enabled: %v)\n", b.Name, *b.Enabled)
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

	config.DebugPrint("DEBUG: Final bead configurations:")
	for _, b := range validBeads {
		fmt.Printf("  - %s (Enabled: %v)\n", b.Name, *b.Enabled)
	}

	config.DebugPrint("DEBUG: Valid beads:")
	for _, b := range validBeads {
		fmt.Printf("  - %s (Enabled: %v)\n", b.Name, *b.Enabled)
	}
	
	processed := make(map[string]int)
	
	for _, b := range validBeads {
		config.DebugPrint("DEBUG: Main loop processing bead %s (Enabled: %v)\n", b.Name, *b.Enabled)
		if b.Enabled != nil && !*b.Enabled {
			config.DebugPrint("DEBUG: Skipping disabled bead in main loop: %s\n", b.Name)
			continue
		}
		if err := processBead(b, yamlData, beadMap, processed, &processedBeads, applyPlan, "", false); err != nil {
			log.Fatalf("Failed to process bead %s: %v", b.Name, err)
		}
	}

	for beadIndex, reason := range invalidBeadReasons {
		beadName := fmt.Sprintf("bead_%s", beadIndex)
		fmt.Printf("Skipping bead: %s, Reason: %s\n", beadName, reason)
		invalidBeadNames = append(invalidBeadNames, fmt.Sprintf("%s: %s", beadName, reason))
	}	

	fmt.Println("\nDEBUG: Processed beads:")
	for _, name := range processedBeads {
		fmt.Printf("  - %s\n", name)
	}
	
	fmt.Println("\nDEBUG: Skipped beads:")
	for name, reason := range invalidBeadReasons {
		fmt.Printf("  - %s: %s\n", name, reason)
	}

}

func handleConfigCommand() {
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
}

func handleFormatCommand() {
	dir := "."
	if len(os.Args) > 2 {
		dir = os.Args[2]
	}
	err := engine.FormatKDFilesInDir(dir)
	if err != nil {
		log.Fatalf("Error formatting .kd files: %v", err)
	}
}