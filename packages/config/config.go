package config

import (
	"bufio"
	"fmt"
	"github.com/janpreet/kado/packages/bead"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type YAMLConfig map[string]interface{}

var LandingZone = "LandingZone"
var TemplateDir = "templates"
var Debug bool = false

const Version = "1.0.0"

func LoadBeadsConfig(filename string) ([]bead.Bead, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var beads []bead.Bead
	scanner := bufio.NewScanner(file)
	var currentBead *bead.Bead
	for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        if strings.HasPrefix(line, "bead \"") {
            if currentBead != nil {
                beads = append(beads, *currentBead)
            }
            currentBead = &bead.Bead{
                Fields: make(map[string]string),
            }
            currentBead.Name = strings.Trim(line[6:], "\" {")
            DebugPrint("DEBUG: Loading bead: %s\n", currentBead.Name)
        } else if currentBead != nil {
            parts := strings.SplitN(line, "=", 2)
            if len(parts) != 2 {
                continue
            }
            key := strings.TrimSpace(parts[0])
            value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
			if strings.ToLower(key) == "enabled" {
				enabled := strings.ToLower(value) == "true"
				currentBead.Enabled = &enabled
				DebugPrint("DEBUG: Set %s.Enabled = %v\n", currentBead.Name, *currentBead.Enabled)
			} else {
				currentBead.Fields[key] = value
			}
        }
    }
    if currentBead != nil {
        beads = append(beads, *currentBead)
    }
	if err := scanner.Err(); err != nil {
		return nil, err
	}

    for i, b := range beads {
        DebugPrint("DEBUG: Loaded bead %s (index: %d) with enabled = %v\n", b.Name, i, b.Enabled)
    }

    return beads, nil
}

func LoadYAMLConfig(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func GetValidBeadNames() map[string]struct{} {
	return GetValidBeads()
}

func GetValidBeadsWithDefaultEnabled(beads []bead.Bead) ([]bead.Bead, map[string]string) {
    var validBeads []bead.Bead
    invalidBeadReasons := make(map[string]string)

    for _, b := range beads {
        if b.Enabled == nil {
            defaultEnabled := true
            b.Enabled = &defaultEnabled
        }
        
        DebugPrint("DEBUG: Validating bead %s (Enabled: %v)\n", b.Name, *b.Enabled)
        
        if *b.Enabled {
            validBeads = append(validBeads, b)
        } else {
            invalidBeadReasons[b.Name] = fmt.Sprintf("%s (disabled)", b.Name)
        }
    }

    return validBeads, invalidBeadReasons
}

func DebugPrint(format string, a ...interface{}) {
    if Debug {
        fmt.Printf(format, a...)
    }
}