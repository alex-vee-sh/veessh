veessh - Console connection manager (SSH/SFTP/Telnet)

veessh is a Go-based CLI to manage console connection profiles and credentials for
SSH, SFTP, Telnet, and other tools. It orchestrates native clients (ssh, sftp,
telnet) and stores credentials securely with the system keychain.

Installation

1. Ensure Go 1.22+
2. Build:

```bash
go build -o veessh ./cmd/veessh
```

Quick start

```bash
# Add an SSH profile
./veessh add mybox --type ssh --host example.com --user alice \
  --port 22 --identity ~/.ssh/id_ed25519

# List profiles
./veessh list

# Connect
./veessh connect mybox
```

Key features

- Interactive picking with fuzzy search (built-in; optional fzf fallback)
- Favorites and recents; usage tracking updates on successful connect
- ProxyJump support; tag support; JSON output for list/show
- Port-forward presets on profiles (local/remote/dynamic)
- Import/export profiles (YAML), and import from OpenSSH config
- Shell completions for bash/zsh/fish/powershell
- Graceful Ctrl+C: clean cancellation with "ok. exiting"

Notes

- veessh uses your system's native tools (ssh, sftp, telnet). Ensure they are
  installed and in PATH.
- Passwords are optional; when provided, they are stored in the OS keychain via
  github.com/99designs/keyring. The current version does not auto-inject
  passwords into ssh; you'll be prompted by the tool as usual. Keys and agents
  are fully supported.
- Config lives at ~/.config/sshmgr/config.yaml by default.

Core commands

- add: Create or update a profile.
- list: Show profiles (supports --tag and --json).
- show: Show details for a profile (supports --json).
- connect: Connect using a profile.
- pick: Interactively pick and connect (supports --fzf, --favorites, --tag,
  --recent-first, --print).
- favorite: Toggle favorite flag.
- export / import: Export/import profiles (YAML; no passwords).
- import-ssh: Import from ~/.ssh/config.
- completion: Emit shell completion script.
- remove: Delete a profile (and optionally its stored password).

Examples

SFTP and Telnet:

```bash
./veessh add filesvc --type sftp --host files.example --user alice
./veessh connect filesvc

./veessh add legacy --type telnet --host legacy.example --port 23
./veessh connect legacy
```

Picker, favorites, tags:

```bash
./veessh favorite mybox
./veessh pick --favorites --tag prod --recent-first --fzf
```

JSON output and tags on list/show:

```bash
./veessh list --tag prod --json
./veessh show mybox --json
```

Import/export and OpenSSH import:

```bash
./veessh export --file profiles.yaml
./veessh import --file profiles.yaml --overwrite
./veessh import-ssh --file ~/.ssh/config --group imported --prefix ssh-
```

Completions:

```bash
# Bash
./veessh completion bash > /usr/local/etc/bash_completion.d/veessh

# Zsh
./veessh completion zsh > "${fpath[1]}/_veessh"

# Fish
./veessh completion fish > ~/.config/fish/completions/veessh.fish

# PowerShell
./veessh completion powershell > veessh.ps1
```

Roadmap

- On-connect automation: per-profile commands (tmux attach/new, cd, env vars),
  remote working dir
- Port-forward presets UX: toggle at connect time; named presets per profile
- Profile templates/inheritance and shared read-only org directory merge
- Rich TUI picker with columns and quick actions (connect, edit, copy, favorite)
- Doctor diagnostics: client availability, key perms, agent status, reachability
- Host key trust: strict/lenient modes, first-connect fingerprint verify, pinning
- Secrets integration: 1Password/Bitwarden/AWS Secrets Manager fetch
- Additional transports: AWS SSM, GCP gcloud SSH, mosh, serial, RDP stubs
- Proxy support: SOCKS/HTTP, ProxyCommand, multi-hop chains with saved hops
- Teams/audit: connection audit log (timestamp/duration, no secrets), reports
- Scripting: JSON for all commands, stable schema, machine-mode to avoid prompts
- Packaging: Homebrew tap, release artifacts, static builds, CI lint/test/build
- Sessions: open multiple profiles into tmux windows, named sessions, layouts

Pluggability

Connectors are pluggable; you can add new protocol handlers in
internal/connectors.

License

Apache-2.0


