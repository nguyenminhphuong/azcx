# azcx

[![Release](https://img.shields.io/github/v/release/nguyenminhphuong/azcx?style=flat-square)](https://github.com/nguyenminhphuong/azcx/releases)
[![License](https://img.shields.io/github/license/nguyenminhphuong/azcx?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/nguyenminhphuong/azcx?style=flat-square)](https://goreportcard.com/report/github.com/nguyenminhphuong/azcx)

**azcx** is a tool to switch between Azure subscriptions faster.

If you manage multiple Azure subscriptions (dev, staging, prod, client projects...), you know how slow `az account set` can be. **azcx** makes switching instant.

### Interactive Mode

Just run `azcx` without arguments to get a fuzzy-searchable list:

![azcx interactive](img/interactive.gif)

---

## What It Does

![azcx demo](img/demo.gif)

```
USAGE:
  azcx                       : interactive fuzzy search
  azcx <name>                : switch to subscription by name (or ID)
  azcx -                     : switch to previous subscription
  azcx -c, --current         : show current subscription
  azcx -l, --list            : list all subscriptions
  azcx alias <alias>=<name>  : create an alias
  azcx alias <alias>         : create an alias (interactive)
  azcx alias -d <alias>      : delete an alias
  azcx alias                 : list aliases
  azcx refresh               : refresh subscription list from Azure
  azcx completion <shell>    : generate shell completion
```

### Quick Switch

```bash
# Switch by typing part of the name
$ azcx prod
Switched to my-production-subscription

# Switch back to previous
$ azcx -
Switched to my-dev-subscription
```

### Aliases

Don't want to type long subscription names? Create aliases:

```bash
# Interactive: select from fuzzy finder
$ azcx alias dev
Select subscription for alias 'dev':
> my-long-development-subscription-name
Switched to my-long-development-subscription-name
Created alias 'dev' -> 'my-long-development-subscription-name'

# Or explicit
$ azcx alias prod=my-production-subscription

# Now just use the alias
$ azcx dev
$ azcx prod
```

---

## Why azcx?

| | `az account set` | `azcx` |
|---|---|---|
| **Startup time** | ~500ms (Python) | ~10ms (Go) |
| **API calls** | Yes | No |
| **Interactive selection** | No | Yes (fuzzy finder) |
| **Previous subscription** | Manual | `azcx -` |
| **Aliases** | No | Yes |

**azcx** reads and writes directly to `~/.azure/azureProfile.json` — the same file Azure CLI uses. No API calls, no authentication overhead.

---

## Installation

### Homebrew (macOS/Linux)

```bash
brew install nguyenminhphuong/tap/azcx
```

### Scoop (Windows)

```powershell
scoop bucket add azcx https://github.com/nguyenminhphuong/scoop-bucket
scoop install azcx
```

### Go Install

```bash
go install github.com/nguyenminhphuong/azcx@latest
```

### Download Binary

Download the latest binary from [GitHub Releases](https://github.com/nguyenminhphuong/azcx/releases):

| Platform | Architecture | Download |
|----------|--------------|----------|
| macOS | Apple Silicon (M1/M2) | [azcx-darwin-arm64](https://github.com/nguyenminhphuong/azcx/releases/latest) |
| macOS | Intel | [azcx-darwin-amd64](https://github.com/nguyenminhphuong/azcx/releases/latest) |
| Linux | x64 | [azcx-linux-amd64](https://github.com/nguyenminhphuong/azcx/releases/latest) |
| Linux | ARM64 | [azcx-linux-arm64](https://github.com/nguyenminhphuong/azcx/releases/latest) |
| Windows | x64 | [azcx-windows-amd64.exe](https://github.com/nguyenminhphuong/azcx/releases/latest) |

### Build from Source

```bash
git clone https://github.com/nguyenminhphuong/azcx.git
cd azcx
go build -o azcx .
sudo mv azcx /usr/local/bin/
```

---

## Shell Completion

### Bash

```bash
# Linux
azcx completion bash | sudo tee /etc/bash_completion.d/azcx

# macOS
azcx completion bash > $(brew --prefix)/etc/bash_completion.d/azcx
```

### Zsh

```bash
# Make sure ~/.zsh/completions is in your fpath
azcx completion zsh > ~/.zsh/completions/_azcx
```

### Fish

```bash
azcx completion fish > ~/.config/fish/completions/azcx.fish
```

### PowerShell

```powershell
azcx completion powershell | Out-String | Invoke-Expression

# To load on startup, add to your profile:
# azcx completion powershell | Out-String | Invoke-Expression
```

---

## Configuration

### Config Location

| Platform | Path |
|----------|------|
| macOS/Linux | `~/.config/azcx/config.json` |
| Windows | `%APPDATA%\azcx\config.json` |

### Config File

```json
{
  "previousSubscription": "my-dev-subscription",
  "aliases": {
    "dev": "my-development-subscription",
    "prod": "my-production-subscription",
    "staging": "my-staging-subscription"
  }
}
```

---

## Requirements

- Azure CLI must be installed and logged in (`az login`)
- The tool reads from `~/.azure/azureProfile.json`

---

## Related Tools

- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/) — Official Azure command-line tool

---

## Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

[MIT](LICENSE) — Phuong Nguyen
