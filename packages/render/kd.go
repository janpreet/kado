package render

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/janpreet/kado/packages/bead"
	"github.com/janpreet/kado/packages/config"
)

func ProcessKdFiles(files []string) (map[string]bead.Bead, []string, error) {
	validBeads := config.GetValidBeads()
	kdBeads := make(map[string]bead.Bead)
	var invalidKdBeadNames []string

	for _, file := range files {
		beads, invalidBeadNames, err := parseKdFile(file)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse kd file %s: %v", file, err)
		}
		invalidKdBeadNames = append(invalidKdBeadNames, invalidBeadNames...)

		for _, b := range beads {

			if enabled, ok := b.Fields["enabled"]; ok && enabled == "false" {
				fmt.Printf("Skipping bead %s because it is disabled\n", b.Name)
				continue
			}
			if _, ok := validBeads[b.Name]; !ok {
				invalidKdBeadNames = append(invalidKdBeadNames, b.Name)
			} else {
				kdBeads[b.Name] = b
			}
		}
	}

	return kdBeads, invalidKdBeadNames, nil
}

func parseKdFile(filePath string) ([]bead.Bead, []string, error) {
	var beads []bead.Bead
	var invalidBeadNames []string

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentBead bead.Bead
	var inBead bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if strings.HasPrefix(line, "bead ") {
			if inBead {
				beads = append(beads, currentBead)
			}
			inBead = true
			currentBead = bead.Bead{Fields: make(map[string]string)}
			matches := regexp.MustCompile(`bead "([^"]+)"`).FindStringSubmatch(line)
			if len(matches) > 1 {
				currentBead.Name = matches[1]
			} else {
				invalidBeadNames = append(invalidBeadNames, line)
			}
		} else if inBead {

			if line == "}" {
				beads = append(beads, currentBead)
				inBead = false
			} else {

				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
					currentBead.Fields[key] = value
				}
			}
		}
	}

	if inBead {
		beads = append(beads, currentBead)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	return beads, invalidBeadNames, nil
}

func GetKDFiles(dir string) ([]string, error) {
	var kdFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".kd" {
			kdFiles = append(kdFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return kdFiles, nil
}
