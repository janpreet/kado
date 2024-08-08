package render

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"regexp"
	"github.com/janpreet/kado/packages/config"
	"github.com/janpreet/kado/packages/keybase"	
)

var keybaseNoteRegex = regexp.MustCompile(`{{keybase:note:([^}]+)}}`)

func join(data map[string]interface{}, key, delimiter string) string {
	var result []string
	for i := 0; ; i++ {
		arrayKey := fmt.Sprintf("%s[%d]", key, i)
		if val, ok := data[arrayKey]; ok {
			result = append(result, fmt.Sprintf("%v", val))
		} else {
			break
		}
	}
	return strings.Join(result, delimiter)
}

type FlattenedDataMap struct {
	Data map[string]interface{}
}

func (f FlattenedDataMap) Get(key string) interface{} {
	if val, ok := f.Data[key]; ok {
		return val
	}
	return "<no value>"
}

func (f FlattenedDataMap) Env(key string) string {
	return os.Getenv(key)
}

func (f FlattenedDataMap) GetKeysAsArray(key string) string {
	keys := strings.Split(key, ".")
	var nestedMap map[string]interface{}

	for k, v := range f.Data {

		if strings.HasPrefix(k, key) && strings.Count(k, ".") == len(keys) {
			if nestedMap == nil {
				nestedMap = make(map[string]interface{})
			}
			parts := strings.Split(k, ".")
			lastPart := parts[len(parts)-1]
			nestedMap[lastPart] = v
		}
	}

	if nestedMap != nil {
		keysArray := make([]string, 0, len(nestedMap))
		for k := range nestedMap {

			strippedKey := strings.Split(k, "[")[0]
			if !contains(keysArray, strippedKey) {
				keysArray = append(keysArray, fmt.Sprintf("\"%s\"", strippedKey))
			}
		}
		return fmt.Sprintf("[%s]", strings.Join(keysArray, ", "))
	}

	return "[]"
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ProcessTemplate(templatePath string, data map[string]interface{}) (string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			filteredLines = append(filteredLines, line)
			continue
		}
		filteredLines = append(filteredLines, line)
	}

	firstLine, templateContent := filteredLines[0], strings.Join(filteredLines[1:], "\n")
	if !strings.HasPrefix(firstLine, "<") || !strings.HasSuffix(firstLine, ">") {
		return "", fmt.Errorf("invalid file name format in template: %s", firstLine)
	}
	fileName := strings.Trim(firstLine, "<>")

	flatData := FlattenYAML("", data)

	funcMap := template.FuncMap{
		"join": func(key, delimiter string) string {
			return join(flatData, key, delimiter)
		},
		"Get": func(key string) interface{} {
			return FlattenedDataMap{Data: flatData}.Get(key)
		},
		"Env": func(key string) string {
			return FlattenedDataMap{Data: flatData}.Env(key)
		},
		"GetKeysAsArray": func(key string) string {
			return FlattenedDataMap{Data: flatData}.GetKeysAsArray(key)
		},
        "KeybaseNote": func(noteName string) (string, error) {
            return resolveKeybaseNote(noteName)
        },		
	}

	processedContent := keybaseNoteRegex.ReplaceAllStringFunc(templateContent, func(match string) string {
        noteName := strings.TrimPrefix(strings.TrimSuffix(match, "}}"), "{{keybase:note:")
        return fmt.Sprintf(`{{ KeybaseNote "%s" }}`, noteName)
    })	

    tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).Parse(processedContent)
    if err != nil {
        return "", fmt.Errorf("failed to parse template: %v", err)
    }

	var output bytes.Buffer
	if err := tmpl.Execute(&output, FlattenedDataMap{Data: flatData}); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	outputPath := filepath.Join(config.LandingZone, fileName)
	err = WriteToFile(outputPath, output.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to write output file: %v", err)
	}

	return outputPath, nil
}

func ProcessTemplates(templatePaths []string, data map[string]interface{}) error {
	for _, templatePath := range templatePaths {
		if _, err := ProcessTemplate(templatePath, data); err != nil {
			return fmt.Errorf("failed to process template %s: %v", templatePath, err)
		}
	}
	return nil
}

func resolveKeybaseNote(noteName string) (string, error) {
    content, err := keybase.ViewNote(noteName)
    if err != nil {
        return "", fmt.Errorf("failed to resolve Keybase note %s: %v", noteName, err)
    }
    return strings.TrimSpace(content), nil
}