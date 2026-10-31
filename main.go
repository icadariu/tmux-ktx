package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/icadariu/tmux-ktx/internal/kube"
)

func formatTmuxOutput(ctx, ns, ctxColor, nsColor string) string {
	// tmux expands #(...) output as a format; "##" keeps names from a kubeconfig literal
	ctx, ns = strings.ReplaceAll(ctx, "#", "##"), strings.ReplaceAll(ns, "#", "##")
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
		_ = os.Setenv("KUBECONFIG", flag.Arg(0))
	}

	ctx, ns, err := kube.CurrentContextAndNamespace()
	if err != nil || ctx == "" {
		os.Exit(0)
	}

	fmt.Print(formatTmuxOutput(ctx, ns, *ctxColor, *nsColor))
}
