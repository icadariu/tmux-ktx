package kube

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type kubeConfig struct {
	CurrentContext string    `yaml:"current-context"`
	Contexts       []context `yaml:"contexts"`
}

type context struct {
	Name    string         `yaml:"name"`
	Context contextDetails `yaml:"context"`
}

type contextDetails struct {
	Namespace string `yaml:"namespace"`
}

// CurrentContextAndNamespace returns the active Kubernetes context name and
// namespace. It respects KUBECONFIG (colon-separated paths) and falls back to
// ~/.kube/config. Returns empty strings without an error when no context is
// configured.
func CurrentContextAndNamespace() (string, string, error) {
	paths := configPaths()

	var currentContext string
	contexts := make(map[string]contextDetails)

	for _, p := range paths {
		cfg, err := parseFile(p)
		if err != nil {
			// skip unreadable or missing files gracefully
			continue
		}

		if currentContext == "" && cfg.CurrentContext != "" {
			currentContext = cfg.CurrentContext
		}

		for _, ctx := range cfg.Contexts {
			if _, exists := contexts[ctx.Name]; !exists {
				contexts[ctx.Name] = ctx.Context
			}
		}
	}

	if currentContext == "" {
		return "", "", nil
	}

	details, ok := contexts[currentContext]
	if !ok {
		return currentContext, "default", nil
	}

	ns := details.Namespace
	if ns == "" {
		ns = "default"
	}

	return currentContext, ns, nil
}

func configPaths() []string {
	if env := os.Getenv("KUBECONFIG"); env != "" {
		return strings.Split(env, ":")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	return []string{filepath.Join(home, ".kube", "config")}
}

func parseFile(path string) (*kubeConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg kubeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
