# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

`make` with no arguments prints the self-documented target list (see Makefile for full set). Key targets:

- `make build` — compile `./tmux-ktx` with version metadata baked in via `-ldflags`.
- `make install` — `go install` into `$GOPATH/bin` with the same ldflags. This is how the user runs after `git pull` to refresh the deployed binary.
- `make test` — runs `go test ./...`.
- `make version` — prints the version string that *would* be embedded (sanity check before building).

Run a single test case:

```bash
go test ./internal/kube -run 'TestCurrentContextAndNamespace/namespace_omitted_defaults_to_default' -v
```

Build-time vars (`VERSION`, `COMMIT`, `BUILDTIME`) are derived from `git describe --tags --always --dirty`, `git rev-parse --short HEAD`, and `date +%y-%m-%d_%H:%M`. They can be overridden on the make line.

## Architecture

This is a single-binary tmux plugin invoked **per-pane on every status refresh** (default 5s). Two design constraints follow:

1. **Must be fast and silent.** Any non-empty stderr or stdout becomes visible noise in the tmux status bar. The error path in `main.go` calls `os.Exit(0)` with no output — that's intentional, not a bug. Do not "improve" it to print errors.
2. **Must not depend on shell environment.** The tmux server has a frozen snapshot of the environment from when the session started; variables `export`ed inside a pane are invisible to it. The plugin solves this by reading `KUBECONFIG` as a **positional argument** passed by tmux (not from `os.Getenv`). The README documents the required zsh `precmd` hook that writes `$KUBECONFIG` into the per-pane tmux option `@ktx_kubeconfig`, which tmux then expands into the argument list.

### Flow

```
zsh precmd hook  →  tmux set-option -p @ktx_kubeconfig "$KUBECONFIG"
                                 ↓
tmux status-right  →  tmux-ktx -ctx-color green -ns-color red #{@ktx_kubeconfig}
                                 ↓
main.go            →  os.Setenv("KUBECONFIG", flag.Arg(0))   ← positional arg → env var
                                 ↓
kube.CurrentContextAndNamespace()  →  merges all KUBECONFIG paths (colon-separated)
                                       first current-context wins
                                       first context definition wins for duplicates
                                 ↓
stdout             →  "#[fg=blue]⎈ #[fg=green]ctx#[fg=colour250]:#[fg=red]ns"
                       (tmux color-escape string, NOT human-readable text)
```

### Package layout

- `main.go` — CLI surface: flags (`-ctx-color`, `-ns-color`, `--version`), KUBECONFIG positional arg, final tmux-formatted printf. The `version`/`commit`/`buildTime` package vars are populated by `-ldflags`.
- `internal/kube/config.go` — kubeconfig parsing. Reads files listed in `KUBECONFIG` (or `~/.kube/config`), merges contexts, returns `(context, namespace, err)`. Missing files and invalid YAML are swallowed (continue) — never errored — because a broken config file in a multi-file `KUBECONFIG` chain must not break the display.
- `internal/kube/config_test.go` — table-driven tests covering merge precedence, missing files, invalid YAML, default namespace fallback, and unknown current-context.

### Defaults that matter

- Namespace falls back to `"default"` when omitted in the kubeconfig OR when `current-context` points to an unknown context.
- Empty `current-context` returns `("", "", nil)` — main.go treats this as "exit silently".
- `KUBECONFIG` unset → falls back to `~/.kube/config`.
