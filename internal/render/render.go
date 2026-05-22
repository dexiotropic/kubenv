package render

import (
	"fmt"
	"regexp"
)

var variablePattern = regexp.MustCompile(`\$([A-Z_][A-Z0-9_]*)`)

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
		key := match[1:]
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
