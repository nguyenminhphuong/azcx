package cmd

import (
	"fmt"
	"strings"

	"github.com/nguyenminhphuong/azcx/internal/config"
	"github.com/spf13/cobra"
)

var deleteAlias string

var aliasCmd = &cobra.Command{
	Use:   "alias [alias=subscription]",
	Short: "Manage subscription aliases",
	Long: `Create or manage short aliases for subscription names.

Examples:
  azcx alias                          # List all aliases
  azcx alias dev=my-dev-subscription  # Create alias 'dev'
  azcx alias -d dev                   # Delete alias 'dev'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			cfg = &config.Config{Aliases: make(map[string]string)}
		}

		// Delete alias
		if deleteAlias != "" {
			if _, ok := cfg.Aliases[deleteAlias]; !ok {
				return fmt.Errorf("alias '%s' not found", deleteAlias)
			}
			delete(cfg.Aliases, deleteAlias)
			if err := config.Save(cfg); err != nil {
				return err
			}
			fmt.Printf("Deleted alias '%s'\n", deleteAlias)
			return nil
		}

		// List aliases
		if len(args) == 0 {
			if len(cfg.Aliases) == 0 {
				fmt.Println("No aliases configured")
				return nil
			}
			for alias, sub := range cfg.Aliases {
				fmt.Printf("%s -> %s\n", alias, sub)
			}
			return nil
		}

		// Create alias
		parts := strings.SplitN(args[0], "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid format. Use: alias=subscription")
		}

		alias, subscription := parts[0], parts[1]
		if alias == "" || subscription == "" {
			return fmt.Errorf("alias and subscription cannot be empty")
		}

		cfg.Aliases[alias] = subscription
		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("Created alias '%s' -> '%s'\n", alias, subscription)
		return nil
	},
}

func init() {
	aliasCmd.Flags().StringVarP(&deleteAlias, "delete", "d", "", "Delete an alias")
	rootCmd.AddCommand(aliasCmd)
}
