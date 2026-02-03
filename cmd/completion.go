package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for azcx.

To load completions:

Bash:
  $ source <(azcx completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ azcx completion bash > /etc/bash_completion.d/azcx
  # macOS:
  $ azcx completion bash > $(brew --prefix)/etc/bash_completion.d/azcx

Zsh:
  $ source <(azcx completion zsh)
  # To load completions for each session, execute once:
  $ azcx completion zsh > "${fpath[1]}/_azcx"

Fish:
  $ azcx completion fish | source
  # To load completions for each session, execute once:
  $ azcx completion fish > ~/.config/fish/completions/azcx.fish

PowerShell:
  PS> azcx completion powershell | Out-String | Invoke-Expression
  # To load completions for each session, execute once:
  PS> azcx completion powershell > azcx.ps1
  # and source this file from your PowerShell profile.
`,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
