package main

import (
	"fmt"
	"os"

	"github.com/icadariu/tmux-ktx/internal/kube"
)

func main() {
	// KUBECONFIG is passed as argv[1], expanded by tmux from the pane option
	// #{@ktx_kubeconfig} that the shell hook keeps up to date.
	// This works for interactive 'export KUBECONFIG=...' unlike reading /proc/environ
	// which only captures the initial process environment.
	if len(os.Args) > 1 && os.Args[1] != "" {
		os.Setenv("KUBECONFIG", os.Args[1])
	}

	ctx, ns, err := kube.CurrentContextAndNamespace()
	if err != nil || ctx == "" {
		os.Exit(0)
	}

	fmt.Printf("#[fg=blue]⎈ #[fg=default]%s#[fg=colour250]:#[fg=default]%s", ctx, ns)
}
