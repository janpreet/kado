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
		} else if currentBead != nil {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
			if key == "enabled" {
				enabled := value == "true"
				currentBead.Enabled = &enabled
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

func GetValidBeadsWithDefaultEnabled(beads []bead.Bead) ([]bead.Bead, []string) {
	var validBeads []bead.Bead
	var invalidBeadReasons []string

	for _, b := range beads {
		if _, ok := ValidBeadNames[b.Name]; !ok {
			invalidBeadReasons = append(invalidBeadReasons, fmt.Sprintf("%s (invalid name)", b.Name))
			continue
		}

		if b.Enabled == nil {
			validBeads = append(validBeads, b)
		} else if !*b.Enabled {
			invalidBeadReasons = append(invalidBeadReasons, fmt.Sprintf("%s (disabled)", b.Name))
		} else {
			validBeads = append(validBeads, b)
		}
	}

	return validBeads, invalidBeadReasons
}
