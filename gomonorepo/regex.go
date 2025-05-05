package gomonorepo

import (
	"regexp"
	"strings"
)

// Patterns is a wrapper around a slice of regular expressions.
type Patterns []*regexp.Regexp

// MatchString checks if any of the patterns match the given string.
func (p Patterns) MatchString(s string) bool {
	for _, pattern := range p {
		if pattern.MatchString(s) {
			return true
		}
	}
	return false
}

func matchesAny(cp Patterns, prefixes []string, s string) bool {
	if cp.MatchString(s) {
		return true
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
