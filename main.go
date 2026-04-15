package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/icadariu/tmux-ktx/internal/kube"
)

func main() {
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

	fmt.Printf("#[fg=blue]⎈ #[fg=%s]%s#[fg=colour250]:#[fg=%s]%s", *ctxColor, ctx, *nsColor, ns)
}
