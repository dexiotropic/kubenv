package render

import (
	"fmt"
	"regexp"
)

// Match variables in the form of {{ env.VAR_NAME }}
var variablePattern = regexp.MustCompile(`\{\{\s*env\.([A-Z_][A-Z0-9_]*)\s*\}\}`)

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

func Strict(input []byte, vars map[string]string) ([]byte, error) {
	missing := map[string]struct{}{}
	output := variablePattern.ReplaceAllStringFunc(string(input), func(match string) string {
		// Get matching group, which is the variable name without the "env." prefix
		submatches := variablePattern.FindStringSubmatch(match)

		if len(submatches) < 2 {
			// This should not happen, but just in case, return the original match
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

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing variables: %v", keys(missing))
	}

	return []byte(output), nil
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
	return out
}
