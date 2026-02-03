# azcx

Fast Azure subscription context switcher, inspired by [kubectx](https://github.com/ahmetb/kubectx).

## Installation

### From source

```bash
go install github.com/nguyenminhphuong/azcx@latest
```

### From releases

Download from [GitHub Releases](https://github.com/nguyenminhphuong/azcx/releases).

### Homebrew (coming soon)

```bash
brew install nguyenminhphuong/tap/azcx
```

## Usage

```bash
# Interactive fuzzy finder
azcx

# Switch to subscription by name
azcx my-subscription

# Switch to previous subscription
azcx -

# Show current subscription
azcx -c

# List all subscriptions
azcx -l

# Refresh subscriptions from Azure
azcx refresh

# Create an alias
azcx alias dev=my-development-subscription
azcx dev  # now works!

# Delete an alias
azcx alias -d dev

# Shell completion
azcx completion bash > /etc/bash_completion.d/azcx
azcx completion zsh > "${fpath[1]}/_azcx"
```

## Why azcx?

The Azure CLI (`az account set`) is slow because:
- Python startup overhead (~300ms)
- SDK initialization
- API calls on every command

`azcx` is fast because it:
- Is a single Go binary (~10ms startup)
- Reads/writes directly to `~/.azure/azureProfile.json`
- No Azure API calls for switching

## License

MIT
