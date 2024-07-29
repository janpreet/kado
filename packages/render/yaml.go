package render

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ProcessYAMLFiles(files []string) ([]map[string]interface{}, error) {
	var parsedYAMLs []map[string]interface{}
	for _, file := range files {
		content, err := ProcessYAMLFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to process YAML file %s: %v", file, err)
		}
		parsedYAMLs = append(parsedYAMLs, content)
	}
	return parsedYAMLs, nil
}

func ProcessYAMLFile(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %v", err)
	}

	var content map[string]interface{}
	if err := yaml.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML content: %v", err)
	}

	flatContent := FlattenYAML("", content)

	return flatContent, nil
}

func FlattenYAML(prefix string, content interface{}) map[string]interface{} {
	flatMap := make(map[string]interface{})

	switch content := content.(type) {
	case map[string]interface{}:
		for key, value := range content {
			flatKey := key
			if prefix != "" {
				flatKey = prefix + "." + key
			}
			for k, v := range FlattenYAML(flatKey, value) {
				flatMap[k] = v
			}
		}
	case []interface{}:
		for i, value := range content {
			flatKey := fmt.Sprintf("%s[%d]", prefix, i)
			flatMap[flatKey] = value
		}
	default:
		flatMap[prefix] = content
	}

	return flatMap
}
