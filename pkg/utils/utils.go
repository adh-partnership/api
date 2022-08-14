package utils

import (
	"regexp"
	"strings"
)

func ArrayContains(array []string, item string) bool {
	for _, a := range array {
		if a == item {
			return true
		}
	}
	return false
}

func StringToSlug(s string) string {
	if len(s) > 100 {
		s = s[:99]
	}
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9 -]+`)
	s = re.ReplaceAllString(s, "")
	s = strings.Replace(s, " ", "-", -1)
	s = strings.Replace(s, "--", "-", -1)
	s = strings.TrimRight(s, "-")
	s = strings.TrimSpace(s)
	return s
}
