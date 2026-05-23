package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "tmux-ktx-test-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "mkdir temp: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmp)

	binaryPath = filepath.Join(tmp, "tmux-ktx-test")
	build := exec.Command("go", "build", "-o", binaryPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func runBinary(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	// isolate from the developer's real KUBECONFIG
	cmd.Env = append(os.Environ(), "KUBECONFIG=")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exit := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exit = exitErr.ExitCode()
	} else if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	return stdout.String(), stderr.String(), exit
}

func TestVersionFlagPrintsCanonicalFormat(t *testing.T) {
	stdout, stderr, exit := runBinary(t, "--version")
	if exit != 0 {
		t.Fatalf("exit %d, stderr=%q", exit, stderr)
	}
	line := strings.TrimRight(stdout, "\n")
	// Expected shape: "<version> (built <YY-MM-DD_HH:MM>, commit <sha>)" (no app name)
	re := regexp.MustCompile(`^\S+ \(built \S+, commit \S+\)$`)
	if !re.MatchString(line) {
		t.Errorf("version line %q does not match canonical format %q", line, re)
	}
}

func TestSilentExitWhenNoContext(t *testing.T) {
	// pass a non-existent kubeconfig as positional arg → no context found
	stdout, stderr, exit := runBinary(t, filepath.Join(t.TempDir(), "does-not-exist"))
	if exit != 0 {
		t.Errorf("expected exit 0, got %d", exit)
	}
	if stdout != "" {
		t.Errorf("expected empty stdout (would pollute tmux status bar), got %q", stdout)
	}
	if stderr != "" {
		t.Errorf("expected empty stderr, got %q", stderr)
	}
}

func TestTmuxFormattedOutputFromKubeconfig(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "config")
	content := `
apiVersion: v1
kind: Config
current-context: kind-kind
contexts:
- name: kind-kind
  context:
    namespace: my-ns
`
	if err := os.WriteFile(cfgPath, []byte(content), 0600); err != nil {
		t.Fatalf("write kubeconfig: %v", err)
	}

	stdout, stderr, exit := runBinary(t, "-ctx-color", "green", "-ns-color", "red", cfgPath)
	if exit != 0 {
		t.Fatalf("exit %d, stderr=%q", exit, stderr)
	}
	want := "#[fg=blue]⎈ #[fg=green]kind-kind#[fg=colour250]:#[fg=red]my-ns"
	if stdout != want {
		t.Errorf("tmux output:\n  got  %q\n  want %q", stdout, want)
	}
}

func TestDefaultColorsWhenFlagsOmitted(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "config")
	content := `
apiVersion: v1
kind: Config
current-context: prod
contexts:
- name: prod
  context: {}
`
	if err := os.WriteFile(cfgPath, []byte(content), 0600); err != nil {
		t.Fatalf("write kubeconfig: %v", err)
	}

	stdout, _, exit := runBinary(t, cfgPath)
	if exit != 0 {
		t.Fatalf("exit %d", exit)
	}
	want := "#[fg=blue]⎈ #[fg=default]prod#[fg=colour250]:#[fg=default]default"
	if stdout != want {
		t.Errorf("default-color output:\n  got  %q\n  want %q", stdout, want)
	}
}
