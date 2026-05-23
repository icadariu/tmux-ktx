package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/icadariu/tmux-ktx/internal/kube"
)

var (
	version   = "dev"
	commit    = "none"
	buildTime = "unknown"
)

// HandleVersionFlag prints the version line and returns true if --version
// (or -version) was passed. Call at the top of main():
//
//	if HandleVersionFlag() { return }
func HandleVersionFlag() bool {
	for _, a := range os.Args[1:] {
		if a == "--version" || a == "-version" {
			fmt.Printf("%s (built %s, commit %s)\n", version, buildTime, commit)
			return true
		}
	}
	return false
}

func formatTmuxOutput(ctx, ns, ctxColor, nsColor string) string {
	return fmt.Sprintf("#[fg=blue]⎈ #[fg=%s]%s#[fg=colour250]:#[fg=%s]%s", ctxColor, ctx, nsColor, ns)
}

func main() {
	if HandleVersionFlag() {
		return
	}

	ctxColor := flag.String("ctx-color", "default", "tmux color for the kubernetes context")
	nsColor := flag.String("ns-color", "default", "tmux color for the kubernetes namespace")
	flag.Parse()

	// KUBECONFIG can be passed as a positional argument after flags
	if flag.NArg() > 0 && flag.Arg(0) != "" {
		os.Setenv("KUBECONFIG", flag.Arg(0))
	}

	ctx, ns, err := kube.CurrentContextAndNamespace()
	if err != nil || ctx == "" {
		os.Exit(0)
	}

	fmt.Print(formatTmuxOutput(ctx, ns, *ctxColor, *nsColor))
}
