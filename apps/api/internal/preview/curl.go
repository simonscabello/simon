package preview

import (
	"strings"
)

func ShellSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, `'`, `'\''`) + "'"
}

func methodAllowsBody(method string) bool {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case "POST", "PUT":
		return true
	default:
		return false
	}
}

type HeaderLine struct {
	Key   string
	Value string
}

func BuildCurl(method, resolvedURL string, headers []HeaderLine, resolvedBody string) string {
	m := strings.ToUpper(strings.TrimSpace(method))
	parts := []string{
		"curl",
		"-X", ShellSingleQuote(m),
		ShellSingleQuote(resolvedURL),
	}
	for _, h := range headers {
		parts = append(parts, "-H", ShellSingleQuote(h.Key+": "+h.Value))
	}
	if methodAllowsBody(m) && resolvedBody != "" {
		parts = append(parts, "--data-raw", ShellSingleQuote(resolvedBody))
	}
	return strings.Join(parts, " ")
}
