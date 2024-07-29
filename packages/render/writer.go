package render

import (
	"fmt"
	"github.com/janpreet/kado/packages/config"
	"os"
	"path/filepath"
)

func WriteToFile(filePath string, data []byte) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func WriteExtraVarsFile(parsedYAMLs []map[string]interface{}, format string) (string, error) {
	var fileName string
	switch format {
	case "yaml":
		fileName = "extra_vars.yaml"
	case "tfvars":
		fileName = "extra_vars.tfvars"
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	filePath := filepath.Join(config.LandingZone, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create extra vars file: %v", err)
	}
	defer file.Close()

	for _, yamlData := range parsedYAMLs {
		for key, value := range yamlData {
			line := fmt.Sprintf("%s = %v\n", key, value)
			if _, err := file.WriteString(line); err != nil {
				return "", fmt.Errorf("failed to write to extra vars file: %v", err)
			}
		}
	}

	return filePath, nil
}
