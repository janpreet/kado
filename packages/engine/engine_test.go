package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatKDFilesInDir(t *testing.T) {
	err := FormatKDFilesInDir("../../")
	assert.NoError(t, err)
}
