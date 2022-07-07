package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeWikiText(t *testing.T) {
	var result string

	// Normal Text
	result = SanitizeWikiText("abcd efgh")
	assert.Equal(t, "abcd efgh", result)

	// Internal Links
	result = SanitizeWikiText("abcd [[cat]] efgh")
	assert.Equal(t, "abcd cat efgh", result)

	// Internal Links with Text
	result = SanitizeWikiText("abcd [[cat|dog]] efgh")
	assert.Equal(t, "abcd dog efgh", result)

	// Internal Links with Text and Spaces
	result = SanitizeWikiText("abcd[[ cat | dog ]]efgh")
	assert.Equal(t, "abcd dog efgh", result)

	// External Links
	result = SanitizeWikiText("abcd [https://thwiki.cc/] efgh")
	assert.Equal(t, "abcd  efgh", result)

	// External Links with Text
	result = SanitizeWikiText("abcd [https://thwiki.cc/ dog] efgh")
	assert.Equal(t, "abcd dog efgh", result)

	// External Links with Text and Spaces
	result = SanitizeWikiText("abcd[https://thwiki.cc/  dog ]efgh")
	assert.Equal(t, "abcd dog efgh", result)

	// Normal Quotes
	result = SanitizeWikiText("abcd 'cat' efgh")
	assert.Equal(t, "abcd 'cat' efgh", result)

	// Bold
	result = SanitizeWikiText("abcd '''cat''' efgh")
	assert.Equal(t, "abcd <b>cat</b> efgh", result)

	// Bold with Spaces
	result = SanitizeWikiText("abcd''' cat '''efgh")
	assert.Equal(t, "abcd<b> cat </b>efgh", result)

	// Italic
	result = SanitizeWikiText("abcd ''cat'' efgh")
	assert.Equal(t, "abcd <i>cat</i> efgh", result)

	// Italic with Spaces
	result = SanitizeWikiText("abcd'' cat ''efgh")
	assert.Equal(t, "abcd<i> cat </i>efgh", result)

	// Bold and Italic
	result = SanitizeWikiText("abcd '''''cat''''' efgh")
	assert.Equal(t, "abcd <i><b>cat</b></i> efgh", result)

	// Bold and Italic with Spaces
	result = SanitizeWikiText("abcd''''' cat '''''efgh")
	assert.Equal(t, "abcd<i><b> cat </b></i>efgh", result)
}
