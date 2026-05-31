package render

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Style controls which placeholder syntax the renderer accepts.
type Style string

const (
	// StyleExplicit renders placeholders in the form {{ env.NAME }}.
	StyleExplicit Style = "explicit"
	// StyleShell renders placeholders in the form $NAME or ${NAME}.
	StyleShell Style = "shell"
)

// Match variables in the form of {{ env.VAR_NAME }}
var explicitVariablePattern = regexp.MustCompile(`\{\{\s*env\.([A-Z_][A-Z0-9_]*)\s*\}\}`)

// Match escaped variables in the form of {{ !env.VAR_NAME }}
var escapedExplicitVariablePattern = regexp.MustCompile(`\{\{\s*!env\.([A-Z_][A-Z0-9_]*)\s*\}\}`)

// Match variables in the form of $VAR or ${VAR}
var shellVariablePattern = regexp.MustCompile(`\$\{?([a-zA-Z_][a-zA-Z0-9_]*)\}?`)

func FromEnviron(environ []string) map[string]string {
	vars := make(map[string]string, len(environ))
	for _, entry := range environ {
		key, value, ok := splitEnv(entry)
		if !ok {
			continue
		}
		vars[key] = value
	}
	return vars
}

// Strict renders using the default explicit placeholder style.
func Strict(input []byte, vars map[string]string) ([]byte, error) {
	return StrictWithStyle(input, vars, StyleExplicit)
}

// StrictWithStyle renders using the provided placeholder style and fails on missing variables.
func StrictWithStyle(input []byte, vars map[string]string, style Style) ([]byte, error) {
	switch style {
	case StyleExplicit:
		output, missing := renderExplicit(string(input), vars)
		if len(missing) > 0 {
			return nil, fmt.Errorf("missing variables: %v", keys(missing))
		}

		return []byte(output), nil
	case StyleShell:
		output, missing := renderWithPattern(string(input), vars, shellVariablePattern)
		if len(missing) > 0 {
			return nil, fmt.Errorf("missing variables: %v", keys(missing))
		}

		return []byte(output), nil
	default:
		return nil, fmt.Errorf("unknown render style: %s", style)
	}
}

func renderExplicit(input string, vars map[string]string) (string, map[string]struct{}) {
	escaped := map[string]string{}
	protected := escapedExplicitVariablePattern.ReplaceAllStringFunc(input, func(match string) string {
		token := fmt.Sprintf("__KUBENV_ESCAPED_%d__", len(escaped))
		escaped[token] = strings.Replace(match, "!env.", "env.", 1)
		return token
	})

	output, missing := renderWithPattern(protected, vars, explicitVariablePattern)
	for token, literal := range escaped {
		output = strings.ReplaceAll(output, token, literal)
	}

	return output, missing
}

func renderWithPattern(input string, vars map[string]string, pattern *regexp.Regexp) (string, map[string]struct{}) {
	missing := map[string]struct{}{}
	output := pattern.ReplaceAllStringFunc(input, func(match string) string {
		submatches := pattern.FindStringSubmatch(match)

		if len(submatches) < 2 {
			return match
		}

		key := submatches[1]
		value, ok := vars[key]
		if !ok {
			missing[key] = struct{}{}
			return match
		}
		return value
	})
	return output, missing
}

func splitEnv(entry string) (string, string, bool) {
	for i := 0; i < len(entry); i++ {
		if entry[i] == '=' {
			return entry[:i], entry[i+1:], true
		}
	}
	return "", "", false
}

func keys(items map[string]struct{}) []string {
	out := make([]string, 0, len(items))
	for key := range items {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
