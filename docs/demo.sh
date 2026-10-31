#!/usr/bin/env bash
# Record a scripted asciinema demo of tmux-ktx (default: docs/demo.cast).
# A throwaway zsh runs inside a private tmux server (own socket, never the
# user's tmux) configured exactly as the README says; asciinema records
# headless by attaching to it while this script types the keystrokes.
# Needs: asciinema >= 3, tmux, zsh, kubectl; agg for the GIF.
# Usage: docs/demo.sh [out.cast]   (run `make build` first; `make demo` does)
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
out=${1:-$root/docs/demo.cast}
size=100x26
ktx=$root/tmux-ktx

for bin in asciinema tmux zsh kubectl; do
  command -v "$bin" >/dev/null || { echo "demo: $bin not on PATH" >&2; exit 1; }
done
[[ -x $ktx ]] || { echo "demo: $ktx missing, run make build" >&2; exit 1; }

tmp=$(mktemp -d /tmp/demo.XXX)
sock="$tmp/tmux.sock"
cleanup() {
  tmux -S "$sock" kill-server 2>/dev/null || true
  rm -rf "$tmp"
}
trap cleanup EXIT

# --- fixtures: sandboxed HOME with three fake kubeconfigs (no real clusters) ---
mkdir -p "$tmp/.kube" "$tmp/bin" "$tmp/zdot"
ln -s "$(command -v kubectl)" "$tmp/bin/kubectl"
kubeconfig() { # file current-context ctx[:ns]...
  local f=$1 cur=$2 c
  shift 2
  {
    echo "apiVersion: v1"
    echo "kind: Config"
    echo "current-context: $cur"
    echo "clusters:"
    for c in "$@"; do printf -- '- name: %s\n  cluster: {server: "https://127.0.0.1:6443"}\n' "${c%%:*}"; done
    echo "users:"
    echo "- name: demo"
    echo "  user: {}"
    echo "contexts:"
    for c in "$@"; do
      printf -- '- name: %s\n  context: {cluster: %s, user: demo' "${c%%:*}" "${c%%:*}"
      [[ $c == *:* ]] && printf ', namespace: %s' "${c#*:}"
      printf '}\n'
    done
  } >"$tmp/.kube/$f"
}
kubeconfig config kind-dev kind-dev
kubeconfig prod.yaml eks-prod eks-prod:monitoring aks-staging:payments

# zsh with the exact precmd hook from the README.
cat >"$tmp/zdot/.zshrc" <<EOF
export PATH="$tmp/bin:\$PATH"
HISTFILE="$tmp/zsh_history"
setopt interactive_comments
PROMPT='%F{green}~/work%f %# '
_tmux_ktx_hook() {
  [[ -z "\$TMUX" ]] && return
  local current="\${KUBECONFIG:-}"
  [[ "\$current" == "\$_tmux_ktx_last" ]] && return
  _tmux_ktx_last="\$current"
  tmux set-option -p @ktx_kubeconfig "\$current" 2>/dev/null
}
precmd_functions+=(_tmux_ktx_hook)
EOF

# tmux config from the README (status-interval lowered so the demo is snappy).
fmt="#($ktx -ctx-color green -ns-color red #{q:@ktx_kubeconfig})"
cat >"$tmp/tmux.conf" <<EOF
set -g default-terminal tmux-256color
set -g pane-border-status top
set -g pane-border-format "[#{?pane_active,#[bold],}#P - $fmt #[nobold]]"
set -g status-left ""
set -g status-style bg=colour236,fg=colour250
set -g status-right "$fmt "
set -g status-interval 1
EOF

t() { tmux -S "$sock" "$@"; }
# typ: send text one char at a time, like a person.
typ() {
  local s=$1 i
  for ((i = 0; i < ${#s}; i++)); do
    t send-keys -l "${s:i:1}"
    sleep 0.05
  done
}
key() { t send-keys "$@"; }
say() { typ "# $1"; sleep 0.5; key Enter; sleep 0.6; }
run() { typ "$1"; sleep 0.3; key Enter; sleep "${2:-2.5}"; }

shell="env HOME=$tmp ZDOTDIR=$tmp/zdot TERM=tmux-256color zsh -d -i"
# Server env sandboxed too: its #() jobs fall back to $HOME/.kube/config.
env -u KUBECONFIG HOME="$tmp" tmux -S "$sock" -f "$tmp/tmux.conf" new-session -d -x "${size%x*}" -y "${size#*x}" "$shell"
mkdir -p "$(dirname "$out")"
asciinema rec --overwrite --headless --window-size "$size" -i 2 \
  -t "tmux-ktx demo" -c "tmux -S $sock attach" "$out" &
recpid=$!
sleep 2

# --- scenes ---
say "pane border + status bar: context:namespace of ~/.kube/config"
sleep 1.5
say "switch namespace, the border follows"
run "kubectl config set-context --current --namespace=web"
say "point this pane at another kubeconfig"
run "export KUBECONFIG=~/.kube/prod.yaml"
say "new pane keeps the default config: each pane shows its own"
key C-l; sleep 0.3
t split-window -h "$shell"; sleep 1.5
run "kubectl config current-context"
t select-pane -L; sleep 0.8
key C-l; sleep 0.3
say "switch context in the left pane"
run "kubectl config use-context aks-staging" 3

sleep 3
t kill-server
wait "$recpid" || true
# Cut tmux's screen restore + "[server exited]" so the cast ends on the last scene.
sed -i '/?1049l/,$d' "$out"
echo "demo: wrote $out (play: asciinema play $out)" >&2
