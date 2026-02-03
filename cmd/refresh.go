package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh subscription list from Azure",
	Long:  `Fetches the latest subscription list from Azure by running 'az account list --refresh'.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Refreshing subscriptions from Azure...")

		azCmd := exec.Command("az", "account", "list", "--refresh", "--output", "none")
		if err := azCmd.Run(); err != nil {
			return fmt.Errorf("failed to refresh: %w (is Azure CLI installed?)", err)
		}

		fmt.Println("Subscriptions refreshed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
