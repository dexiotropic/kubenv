package cmp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"strconv"
	"strings"

	"github.com/dexiotropic/kubenv/internal/render"
	"github.com/dexiotropic/kubenv/internal/version"
	ucli "github.com/urfave/cli/v3"
)

// Run is the entrypoint for the Argo CD CMP binary.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer, environ []string) error {
	command := &ucli.Command{
		Name:        "kubenv-argocd-cmp",
		Usage:       "Render manifests for Argo CD Config Management Plugins",
		Description: "Reads manifests from stdin and resolves values from ARGOCD_APP_PARAMETERS, ARGOCD_ENV_* and the remaining process environment.",
		Version:     version.String(),
		Reader:      stdin,
		Writer:      stdout,
		ErrWriter:   stderr,
		OnUsageError: func(_ context.Context, _ *ucli.Command, err error, _ bool) error {
			return err
		},
		Action: func(_ context.Context, _ *ucli.Command) error {
			input, err := io.ReadAll(stdin)
			if err != nil {
				return err
			}

			vars, style, err := loadVariables(environ)
			if err != nil {
				return err
			}

			output, err := render.StrictWithStyle(input, vars, style)
			if err != nil {
				return err
			}

			_, err = stdout.Write(output)
			return err
		},
	}

	return command.Run(context.Background(), append([]string{"kubenv-argocd-cmp"}, args...))
}

// loadVariables processes the environment variables according to the precedence rules defined by Argo CD CMP.
// For reference, see https://argo-cd.readthedocs.io/en/stable/operator-manual/config-management-plugins/#using-environment-variables-in-your-plugin
func loadVariables(environ []string) (map[string]string, render.Style, error) {
	// Process environment variables take the lowest precedence.
	vars := render.FromEnviron(environ)
	style := render.StyleExplicit

	// Then we include all environment variables with the ARGOCD_ENV_ prefix, stripping the prefix first.
	for _, entry := range environ {
		key, value, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}

		// Argo CD injects environment variables with the ARGOCD_ENV_ prefix for each parameter defined in the Application manifest.
		if after, found := strings.CutPrefix(key, "ARGOCD_ENV_"); found {
			vars[after] = value
		}
	}

	// Finally, we include the ARGOCD_APP_PARAMETERS, which have the highest precedence.
	parameters, options, err := loadParameters(environ)
	if err != nil {
		return nil, "", err
	}

	maps.Copy(vars, parameters)
	if options.shellStyle {
		style = render.StyleShell
	}

	return vars, style, nil
}

func loadParameters(environ []string) (map[string]string, cmpOptions, error) {
	raw := lookupEnv(environ, "ARGOCD_APP_PARAMETERS")
	if raw == "" {
		return map[string]string{}, cmpOptions{}, nil
	}

	var parameters []parameter
	if err := json.Unmarshal([]byte(raw), &parameters); err != nil {
		return nil, cmpOptions{}, fmt.Errorf("invalid ARGOCD_APP_PARAMETERS: %w", err)
	}

	vars := map[string]string{}
	options := cmpOptions{}
	for _, parameter := range parameters {
		if parameter.Name == "kubenv" {
			if err := options.apply(parameter.Map); err != nil {
				return nil, cmpOptions{}, err
			}
			continue
		}

		if parameter.String != nil {
			vars[parameter.Name] = *parameter.String
		}

		maps.Copy(vars, parameter.Map)

		for index, value := range parameter.Array {
			vars[fmt.Sprintf("%s_%d", parameter.Name, index)] = value
		}
	}

	return vars, options, nil
}

type cmpOptions struct {
	shellStyle bool
}

func (options *cmpOptions) apply(values map[string]string) error {
	for key, value := range values {
		switch key {
		case "shell-style":
			enabled, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid kubenv.shell-style value %q: %w", value, err)
			}
			options.shellStyle = enabled
		default:
			return fmt.Errorf("unsupported kubenv option: %s", key)
		}
	}

	return nil
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
