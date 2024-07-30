package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetYAMLFiles(t *testing.T) {
	files, err := GetYAMLFiles("../../")
	assert.NoError(t, err)
	assert.Contains(t, files, "../../cluster.yaml")
}

func TestGetKdFiles(t *testing.T) {
	files, err := GetKdFiles("../../")
	assert.NoError(t, err)
	assert.Contains(t, files, "../../cluster.kd")
}
