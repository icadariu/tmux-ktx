# tmux-ktx vs alternatives

How tmux-ktx compares with the projects it was inspired by (see
[Credits](../README.md#credits)). Facts about the alternatives come from their
READMEs and source on GitHub as of 31-Oct-2026; re-check before relying on them.

## At a glance

| | tmux-ktx | [kube-tmux](https://github.com/jonmosco/kube-tmux) | [tmux-kubecontext](https://github.com/uesyn/tmux-kubecontext) | [tmux-kube](https://github.com/kr3cj/tmux-kube) | [kube-tmuxp](https://github.com/thecasualcoder/kube-tmuxp) |
| --- | --- | --- | --- | --- | --- |
| What it is | Status/border widget | Status widget | Status widget | Status widget | tmuxp session generator |
| Language | Go binary | Bash | Bash | Bash | Go |
| Needs `kubectl` | No | Yes (or `oc`) | Yes | Yes, plus `jq` | Yes, plus `gcloud`, `tmuxp` |
| Per-pane `KUBECONFIG` | Yes, via zsh hook | No, tmux server env | Yes, via hook | No, tmux server env | Per session, by design |
| Per-pane border display | Yes, documented | No | No | No | n/a |
| Caching | None needed | In-process only, lost on every call | kubeconfig mtime, kept in tmux options | None | n/a |
| On error | Prints nothing | `N/A:N/A` | `[E] message` | Empty | n/a |
| TPM install | No | Yes | No | Yes | No (Homebrew) |
| Customization | 2 colors | Many env vars, symbol, divider, transform functions | Many tmux options | Prod/stage/dev colors | n/a |
| Last push | Active | May 2026 | Apr 2021 | Apr 2021 | Feb 2023 |

## Where tmux-ktx is better

* **Correct per pane, not per server.** (High) tmux runs `#()` commands in the
  server's environment, frozen at server start. kube-tmux and tmux-kube read that
  `KUBECONFIG`, so `export KUBECONFIG=...` in a pane never shows up. tmux-ktx gets
  the pane's value through the `@ktx_kubeconfig` pane option and a positional
  argument, so each pane border shows its own context. tmux-kubecontext solves the
  same problem with a similar hook, but only for the status bar.
* **No `kubectl` per refresh.** (High) tmux-ktx parses the kubeconfig YAML itself.
  On this machine (kubectl v1.37.1, small kubeconfig, 50 runs) one call took
  about 6 ms vs about 80 ms for `kubectl config view --minify`, the command
  kube-tmux runs. kube-tmux has an mtime cache, but it lives in shell variables of
  a process that tmux starts fresh on every refresh, so it never hits. With a
  border per pane this cost is multiplied by the pane count every
  `status-interval`.
* **No runtime dependencies.** (High) One static binary. No `kubectl`, `jq`,
  `stat`, or bash version differences (the alternatives carry GNU/BSD `stat`
  branches for that).
* **Tolerant of broken multi-file `KUBECONFIG`.** (High) A missing or invalid
  file in the chain is skipped and the rest still merge, with kubectl's
  first-wins rules. tmux-kubecontext shows an error when any listed file is
  missing. Note this is more lenient than kubectl itself, see below.
* **Safe output.** (Medium) `#` in context or namespace names is escaped, so a
  name can not inject tmux format strings into the status line. kube-tmux
  concatenates the names into its output unescaped; the other two are not
  verified.
* **Tested and released.** (High) Table-driven Go tests, CI, signed releases
  with checksums. The bash alternatives have no tests in their repos.

## Where tmux-ktx is worse

* **Harder install.** (High) No TPM entrypoint. You download a binary or build
  with Go, write an absolute path into `.tmux.conf` and add a shell hook. kube-tmux
  and tmux-kube are one `set -g @plugin` line.
* **Depends on a shell hook, documented only for zsh.** (High) Without the hook it
  falls back to the server's `KUBECONFIG` or `~/.kube/config`, which is no better
  than the alternatives. bash (`PROMPT_COMMAND`) and fish users have to write
  their own. The pane value updates only when a prompt returns.
* **Almost no customization.** (High) Only the context and namespace colors.
  No symbol or divider toggle, no hiding the namespace, no name shortening (long
  EKS ARN contexts fill the status bar), no prod highlighting as tmux-kube has,
  no `oc`. kube-tmux is far ahead here.
* **Silent failure.** (Medium) Printing nothing keeps the status bar clean, but
  a broken kubeconfig looks the same as "no context". kube-tmux shows `N/A` and
  tmux-kubecontext shows the error.
* **Its own parser, not kubectl's.** (Medium) It reimplements the merge rules
  for `current-context` and `contexts[].context.namespace` only. If kubectl's
  merge behavior changes, or an edge case differs, the display can drift from what
  `kubectl` uses; the `kubectl`-based tools match by construction. One known case:
  with an invalid YAML file in `KUBECONFIG`, kubectl v1.37.1 exits with
  `error loading config file`, while tmux-ktx still shows the context from the
  valid files. The status bar then names a context that `kubectl` refuses to use.
  It also can not show cluster or user names.
* **Young, single maintainer.** (High) v0.1.1, built mostly to learn. kube-tmux
  has about 390 stars and a maintainer who ships changes.

## kube-tmuxp is a different tool

kube-tmuxp does not display anything. It generates one tmuxp session per GKE
cluster, each with its own kubeconfig, so switching cluster means switching
session. It needs `gcloud` and is GKE-only. It can be used together with
tmux-ktx: tmux-ktx would show the context of each generated session.

## Which to pick

* Different `KUBECONFIG` per pane, or context on pane borders: tmux-ktx.
* One kubeconfig, TPM install, rich styling: kube-tmux.
* Prod/stage/dev coloring is the main need: tmux-kube, or kube-tmux with a
  custom function.

## Open uncertainties

* The 6 ms vs 80 ms figures are from one machine with a tiny kubeconfig; large
  kubeconfigs and cold caches will change both numbers.
* The alternatives' behavior is read from source, not run. tmux-kubecontext's
  per-pane hook was checked only for the part that matters: tmux does expand
  `#{pane_id}` in the option name.
* kr3cj/tmux-kube is not marked as a fork, but its README installs
  `sudermanjr/tmux-kube`, which no longer exists on GitHub.
