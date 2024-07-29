package config

import (
	"os"
	"path/filepath"
)

func GetYAMLFiles(dir string) ([]string, error) {
	var yamlFiles []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			yamlFiles = append(yamlFiles, filepath.Join(dir, file.Name()))
		}
	}

	return yamlFiles, nil
}

func GetKdFiles(dir string) ([]string, error) {
	var kdFiles []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".kd" {
			kdFiles = append(kdFiles, filepath.Join(dir, file.Name()))
		}
	}

	return kdFiles, nil
}
