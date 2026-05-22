package main

import "testing"

func TestFormatVersion(t *testing.T) {
	got := formatVersion("tmux-ktx", "v0.0.1-1-g98d2612", "22-05-26_20:49", "98d2612")
	want := "tmux-ktx v0.0.1-1-g98d2612 (built 22-05-26_20:49, commit 98d2612)"
	if got != want {
		t.Errorf("formatVersion:\n  got  %q\n  want %q", got, want)
	}
}

func TestFormatTmuxOutput(t *testing.T) {
	got := formatTmuxOutput("kind-kind", "my-namespace", "green", "red")
	want := "#[fg=blue]⎈ #[fg=green]kind-kind#[fg=colour250]:#[fg=red]my-namespace"
	if got != want {
		t.Errorf("formatTmuxOutput:\n  got  %q\n  want %q", got, want)
	}
}

func TestIsVersionRequest(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{"version subcommand", []string{"tmux-ktx", "version"}, true},
		{"--version long flag", []string{"tmux-ktx", "--version"}, true},
		{"-version short flag (Go flag pkg style)", []string{"tmux-ktx", "-version"}, true},
		{"no args", []string{"tmux-ktx"}, false},
		{"unrelated flag", []string{"tmux-ktx", "-ctx-color", "green"}, false},
		{"kubeconfig positional", []string{"tmux-ktx", "/home/user/.kube/config"}, false},
		{"version-like flag among others", []string{"tmux-ktx", "-ctx-color", "green", "--version"}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isVersionRequest(tc.args); got != tc.want {
				t.Errorf("isVersionRequest(%v) = %v, want %v", tc.args, got, tc.want)
			}
		})
	}
}
