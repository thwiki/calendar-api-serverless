package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thwiki/calendar-api-serverless/utils"
)

func TestPrintError(t *testing.T) {
	var err error
	var text string

	// Standard
	err = &utils.InvalidDateError{}
	text = string(printError(err))
	assert.Equal(t, "{\"error\":\""+err.Error()+"\"}", text)
}
