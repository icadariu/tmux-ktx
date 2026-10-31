package kube

import (
	"os"
	"path/filepath"
	"testing"
)

// writeConfig writes content to a file inside dir and returns the full path.
func writeConfig(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writeConfig: %v", err)
	}
	return path
}

func TestCurrentContextAndNamespace(t *testing.T) {
	tests := []struct {
		name    string
		files   map[string]string // filename -> yaml content
		order   []string          // order of files for KUBECONFIG (by filename)
		wantCtx string
		wantNS  string
	}{
		{
			name: "context and namespace set",
			files: map[string]string{
				"config": `
apiVersion: v1
kind: Config
current-context: kind-kind
contexts:
- name: kind-kind
  context:
    cluster: kind-kind
    namespace: my-namespace
`,
			},
			order:   []string{"config"},
			wantCtx: "kind-kind",
			wantNS:  "my-namespace",
		},
		{
			name: "namespace omitted defaults to default",
			files: map[string]string{
				"config": `
apiVersion: v1
kind: Config
current-context: prod
contexts:
- name: prod
  context:
    cluster: prod
`,
			},
			order:   []string{"config"},
			wantCtx: "prod",
			wantNS:  "default",
		},
		{
			name: "no current-context returns empty",
			files: map[string]string{
				"config": `
apiVersion: v1
kind: Config
contexts:
- name: prod
  context:
    cluster: prod
`,
			},
			order:   []string{"config"},
			wantCtx: "",
			wantNS:  "",
		},
		{
			name:    "non-existent file returns empty without error",
			files:   map[string]string{},
			order:   []string{"does-not-exist.yaml"},
			wantCtx: "",
			wantNS:  "",
		},
		{
			name: "invalid yaml returns empty without error",
			files: map[string]string{
				"config": `this: is: not: valid: yaml: [[[`,
			},
			order:   []string{"config"},
			wantCtx: "",
			wantNS:  "",
		},
		{
			name: "current-context points to unknown context returns default namespace",
			files: map[string]string{
				"config": `
apiVersion: v1
kind: Config
current-context: ghost
contexts:
- name: prod
  context:
    cluster: prod
`,
			},
			order:   []string{"config"},
			wantCtx: "ghost",
			wantNS:  "default",
		},
		{
			name: "multiple files - first current-context wins",
			files: map[string]string{
				"first": `
apiVersion: v1
kind: Config
current-context: ctx-a
contexts:
- name: ctx-a
  context:
    namespace: ns-a
`,
				"second": `
apiVersion: v1
kind: Config
current-context: ctx-b
contexts:
- name: ctx-b
  context:
    namespace: ns-b
`,
			},
			order:   []string{"first", "second"},
			wantCtx: "ctx-a",
			wantNS:  "ns-a",
		},
		{
			name: "multiple files - context defined in second file",
			files: map[string]string{
				"first": `
apiVersion: v1
kind: Config
current-context: ctx-b
`,
				"second": `
apiVersion: v1
kind: Config
contexts:
- name: ctx-b
  context:
    namespace: ns-b
`,
			},
			order:   []string{"first", "second"},
			wantCtx: "ctx-b",
			wantNS:  "ns-b",
		},
		{
			name: "multiple files - first file wins for duplicate context",
			files: map[string]string{
				"first": `
apiVersion: v1
kind: Config
current-context: shared
contexts:
- name: shared
  context:
    namespace: from-first
`,
				"second": `
apiVersion: v1
kind: Config
contexts:
- name: shared
  context:
    namespace: from-second
`,
			},
			order:   []string{"first", "second"},
			wantCtx: "shared",
			wantNS:  "from-first",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()

			// Write all config files.
			paths := make(map[string]string)
			for name, content := range tc.files {
				paths[name] = writeConfig(t, dir, name, content)
			}

			// Build the KUBECONFIG value. For non-existent files keep the raw name.
			kubeconfig := ""
			for i, name := range tc.order {
				p, ok := paths[name]
				if !ok {
					p = filepath.Join(dir, name)
				}
				if i == 0 {
					kubeconfig = p
				} else {
					kubeconfig += ":" + p
				}
			}

			t.Setenv("KUBECONFIG", kubeconfig)

			gotCtx, gotNS, err := CurrentContextAndNamespace()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotCtx != tc.wantCtx {
				t.Errorf("context: got %q, want %q", gotCtx, tc.wantCtx)
			}
			if gotNS != tc.wantNS {
				t.Errorf("namespace: got %q, want %q", gotNS, tc.wantNS)
			}
		})
	}
}
