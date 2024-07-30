package bead

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeadCreation(t *testing.T) {
	bead := Bead{
		Name:   "test_bead",
		Fields: map[string]string{"key": "value"},
	}

	assert.Equal(t, "test_bead", bead.Name)
	assert.Equal(t, "value", bead.Fields["key"])
}
