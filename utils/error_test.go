package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintError(t *testing.T) {
	var err error
	var text string

	// Standard
	err = &InvalidDateError{}
	text = string(PrintError(err))
	assert.Equal(t, "{\"error\":\""+err.Error()+"\"}", text)
}
