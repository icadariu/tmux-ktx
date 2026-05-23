package main

import "testing"

func TestFormatTmuxOutput(t *testing.T) {
	got := formatTmuxOutput("kind-kind", "my-namespace", "green", "red")
	want := "#[fg=blue]⎈ #[fg=green]kind-kind#[fg=colour250]:#[fg=red]my-namespace"
	if got != want {
		t.Errorf("formatTmuxOutput:\n  got  %q\n  want %q", got, want)
	}
}
