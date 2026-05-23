package cmp

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/dexiotropic/kubenv/internal/render"
)

// Run is the entrypoint for the Argo CD CMP binary.
func Run(stdin io.Reader, stdout, _ io.Writer, environ []string) error {
	input, err := io.ReadAll(stdin)
	if err != nil {
		return err
	}

	vars, err := loadVariables(environ)
	if err != nil {
		return err
	}

	output, err := render.Strict(input, vars)
	if err != nil {
		return err
	}

	_, err = stdout.Write(output)
	return err
}

func loadVariables(environ []string) (map[string]string, error) {
	vars := render.FromEnviron(environ)

	for _, entry := range environ {
		key, value, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}

		if strings.HasPrefix(key, "ARGOCD_ENV_") {
			vars[strings.TrimPrefix(key, "ARGOCD_ENV_")] = value
		}
	}

	parameters, err := loadParameters(environ)
	if err != nil {
		return nil, err
	}

	for key, value := range parameters {
		vars[key] = value
	}

	return vars, nil
}

func loadParameters(environ []string) (map[string]string, error) {
	raw := lookupEnv(environ, "ARGOCD_APP_PARAMETERS")
	if raw == "" {
		return map[string]string{}, nil
	}

	var parameters []parameter
	if err := json.Unmarshal([]byte(raw), &parameters); err != nil {
		return nil, fmt.Errorf("invalid ARGOCD_APP_PARAMETERS: %w", err)
	}

	vars := map[string]string{}
	for _, parameter := range parameters {
		if parameter.String != nil {
			vars[parameter.Name] = *parameter.String
		}

		for key, value := range parameter.Map {
			vars[key] = value
		}

		for index, value := range parameter.Array {
			vars[fmt.Sprintf("%s_%d", parameter.Name, index)] = value
		}
	}

	return vars, nil
}

type parameter struct {
	Name   string            `json:"name"`
	String *string           `json:"string,omitempty"`
	Map    map[string]string `json:"map,omitempty"`
	Array  []string          `json:"array,omitempty"`
}

func lookupEnv(environ []string, key string) string {
	for _, entry := range environ {
		currentKey, value, ok := strings.Cut(entry, "=")
		if ok && currentKey == key {
			return value
		}
	}

	return ""
}
