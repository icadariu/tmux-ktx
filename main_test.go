package main

import "testing"

func TestFormatTmuxOutput(t *testing.T) {
	got := formatTmuxOutput("kind-kind", "my-namespace", "green", "red")
	want := "#[fg=blue]⎈ #[fg=green]kind-kind#[fg=colour250]:#[fg=red]my-namespace"
	if got != want {
		t.Errorf("formatTmuxOutput:\n  got  %q\n  want %q", got, want)
	}
}

func TestFormatTmuxOutputEscapesHash(t *testing.T) {
	got := formatTmuxOutput("ctx#{pane_id}", "ns#[fg=red]x", "green", "red")
	want := "#[fg=blue]⎈ #[fg=green]ctx##{pane_id}#[fg=colour250]:#[fg=red]ns##[fg=red]x"
	if got != want {
		t.Errorf("formatTmuxOutput:\n  got  %q\n  want %q", got, want)
	}
}
