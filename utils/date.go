package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	dateRe = regexp.MustCompile(`^(\d\d\d\d)[-/](1[012]|0?[1-9])[-/](3[01]|[12][0-9]|0?[1-9])$`)
)

type InvalidDateError struct{}

func (e *InvalidDateError) Error() string {
	return "invalid date"
}

func SanitizeDate(value string) (string, error) {
	matches := dateRe.FindStringSubmatch(value)
	if len(matches) != 4 {
		return "", &InvalidDateError{}
	}
	return fmt.Sprintf("%04s-%02s-%02s", matches[1], matches[2], matches[3]), nil
}

func SanitizeSMWDate(value string) (string, error) {
	parts := strings.Split(value, "/")
	if len(parts) <= 3 {
		return SanitizeDate(value)
	}
	return SanitizeDate(strings.Join(parts[1:4], "-"))
}
