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

func formatVersion(name, version, buildTime, commit string) string {
	return fmt.Sprintf("%s %s (built %s, commit %s)", name, version, buildTime, commit)
}

func printVersion() {
	fmt.Println(formatVersion("tmux-ktx", version, buildTime, commit))
}

func formatTmuxOutput(ctx, ns, ctxColor, nsColor string) string {
	return fmt.Sprintf("#[fg=blue]⎈ #[fg=%s]%s#[fg=colour250]:#[fg=%s]%s", ctxColor, ctx, nsColor, ns)
}

func isVersionRequest(args []string) bool {
	if len(args) < 2 {
		return false
	}
	for _, a := range args[1:] {
		switch a {
		case "version", "-version", "--version":
			return true
		}
	}
	return false
}

func main() {
	if isVersionRequest(os.Args) {
		printVersion()
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
