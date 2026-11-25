package lib

import (
	"regexp"
	"strings"
)

var regex = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = regex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
