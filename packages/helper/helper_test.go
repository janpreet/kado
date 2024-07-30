package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	exists := FileExists("../../cluster.yaml")
	assert.True(t, exists)

	exists = FileExists("../../clusters.yaml")
	assert.False(t, exists)
}
