package gomonorepo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

func matchesString(cp Patterns, p []Patterns, s string) bool {
	if cp.MatchString(s) {
		return true
	}
	for _, pattern := range p {
		if pattern.MatchString(s) {
			return true
		}
	}
	return false
}

// parseGitignore parses a .gitignore content and returns a slice of patterns.
func parseGitignore(filePath string) (patterns Patterns, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	dir := filepath.Dir(filePath)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comments
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		pattern, err := convertGitIgnoreToRegex(dir, line)
		if err != nil {
			return nil, fmt.Errorf("failed to convert pattern %q in file %q: %w", line, filePath, err)
		} else if pattern != nil {
			println(fmt.Sprintf("GAVIN: path: %q, pattern: %q, regex: %q", filePath, line, pattern))
			patterns = append(patterns, pattern)
		}
	}

	return patterns, err
}

// convertGitIgnoreToRegex converts a gitignore pattern to a regular expression.
// Note: is used specifically to find directories that we should not traverse, so it will skip
// some patterns that are not relevant for this use case and return a nil regex when the pattern is skipped.
func convertGitIgnoreToRegex(gitIgnoreDir, pattern string) (*regexp.Regexp, error) {
	// Handle negation pattern;
	if strings.HasPrefix(pattern, "!") {
		return nil, nil
	}

	// Remove trailing slash for directory patterns (caller should handle directory logic)
	pattern = strings.TrimSuffix(pattern, "/")

	if !strings.HasSuffix(pattern, "/") {
		gitIgnoreDir += "/"
	}

	// Start building the regex pattern
	regexPattern := "^" + regexp.QuoteMeta(gitIgnoreDir)

	// If pattern doesn't start with slash, it can match anywhere in the path
	if !strings.HasPrefix(pattern, "/") {
		regexPattern += "(.*/)?"
	} else {
		// Remove the leading slash as this is already included in gitIgnoreDir
		pattern = pattern[1:]
	}

	// Escape special regex characters
	pattern = regexp.QuoteMeta(pattern)

	// Convert gitignore glob patterns to regex
	// Handle double asterisk (match across directories)
	pattern = strings.ReplaceAll(pattern, "\\*\\*", ".*")
	// Handle single asterisk (match within a directory)
	pattern = strings.ReplaceAll(pattern, "\\*", "[^/]*")
	// Handle question mark
	pattern = strings.ReplaceAll(pattern, "\\?", "[^/]")

	regexPattern += pattern + "$"

	return regexp.Compile(regexPattern)
}
