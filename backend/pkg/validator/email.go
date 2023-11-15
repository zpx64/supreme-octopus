package validator

import (
	"strings"
	"unicode"
)

// Accepts at least the x@y.zz pattern.
func IsEmail(v string) bool {
	if v == "" {
		return false
	}
	if containsWhitespace(v) {
		return false
	}

	iAt := strings.IndexByte(v, '@')
	if iAt == -1 {
		return false
	}

	localPart := v[:iAt]
	if localPart == "" {
		return false
	}

	domain := v[iAt+1:]
	if domain == "" {
		return false
	}

	iDot := strings.IndexByte(domain, '.')
	if iDot == -1 || iDot == 0 || iDot == len(domain)-1 {
		return false
	}

	if strings.Index(domain, "..") != -1 {
		return false
	}

	iTLD := strings.LastIndexByte(domain, '.')
	return 2 <= len([]rune(domain[iTLD+1:]))
}

// check if string contains whitespaces
func containsWhitespace(v string) bool {
	for _, r := range v {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}
