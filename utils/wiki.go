package utils

import (
	"regexp"
)

var (
	internalLinksRe = regexp.MustCompile(`\[\[(?:[^\[\|\]]+\|)?([^\[\|\]]+)\]\]`)
	externalLinksRe = regexp.MustCompile(`\[(?:[^ \[\]]+) ?([^\[\]]*)\]`)
	boldRe          = regexp.MustCompile(`'''([^']*)'''`)
	italicRe        = regexp.MustCompile(`''([^']*)''`)
)

func SanitizeWikiText(value string) string {
	value = boldRe.ReplaceAllString(value, "<b>$1</b>")
	value = italicRe.ReplaceAllString(value, "<i>$1</i>")
	value = internalLinksRe.ReplaceAllString(value, "$1")
	value = externalLinksRe.ReplaceAllString(value, "$1")

	return value
}
