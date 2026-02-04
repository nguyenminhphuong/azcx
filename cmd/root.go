package cmd

import (
	"fmt"
	"os"

	"github.com/nguyenminhphuong/azcx/internal/azure"
	"github.com/nguyenminhphuong/azcx/internal/config"
	"github.com/nguyenminhphuong/azcx/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version     = "dev"
	showCurrent bool
	showList    bool
)

func SetVersion(v string) {
	version = v
}

var rootCmd = &cobra.Command{
	Use:   "azcx [subscription]",
	Short: "Fast Azure subscription switcher",
	Long: `azcx - A fast Azure subscription context switcher, inspired by kubectx.

Examples:
  azcx                    # Interactive fuzzy finder
  azcx my-subscription    # Switch to subscription by name
  azcx -                  # Switch to previous subscription
  azcx -c                 # Show current subscription
  azcx -l                 # List all subscriptions`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if showCurrent {
			return runCurrent()
		}
		if showList {
			return runList()
		}
		if len(args) == 1 {
			return runSwitch(args[0])
		}
		return runInteractive()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&showCurrent, "current", "c", false, "Show current subscription")
	rootCmd.Flags().BoolVarP(&showList, "list", "l", false, "List all subscriptions")
	rootCmd.Version = version
}

func runCurrent() error {
	profile, err := azure.LoadProfile()
	if err != nil {
		return err
	}

	current := profile.GetCurrentSubscription()
	if current == nil {
		fmt.Println("No subscription selected")
		return nil
	}

	fmt.Printf("%s (%s)\n", current.GetName(), current.GetID())
	return nil
}

func runList() error {
	profile, err := azure.LoadProfile()
	if err != nil {
		return err
	}

	for _, sub := range profile.Subscriptions {
		marker := "  "
		if sub.IsDefault() {
			marker = "* "
		}
		fmt.Printf("%s%s\n", marker, sub.GetName())
	}
	return nil
}

func runSwitch(target string) error {
	// Handle "-" for previous subscription
	if target == "-" {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.PreviousSubscription == "" {
			return fmt.Errorf("no previous subscription")
		}
		target = cfg.PreviousSubscription
	}

	// Check if target is an alias
	cfg, _ := config.Load()
	if cfg != nil {
		if resolved, ok := cfg.Aliases[target]; ok {
			target = resolved
		}
	}

	profile, err := azure.LoadProfile()
	if err != nil {
		return err
	}

	// Save current as previous before switching
	current := profile.GetCurrentSubscription()
	if current != nil {
		if cfg == nil {
			cfg = &config.Config{Aliases: make(map[string]string)}
		}
		cfg.PreviousSubscription = current.GetName()
		config.Save(cfg)
	}

	// Find and switch to target
	if err := profile.SetSubscription(target); err != nil {
		return err
	}

	if err := azure.SaveProfile(profile); err != nil {
		return err
	}

	fmt.Printf("Switched to %s\n", target)
	return nil
}

func runInteractive() error {
	profile, err := azure.LoadProfile()
	if err != nil {
		return err
	}

	if len(profile.Subscriptions) == 0 {
		return fmt.Errorf("no subscriptions found. Run 'az login' first")
	}

	names := make([]string, len(profile.Subscriptions))
	for i, sub := range profile.Subscriptions {
		names[i] = sub.GetName()
	}

	selected, err := ui.FuzzySelect(names)
	if err != nil {
		return err
	}

	if selected != "" {
		return runSwitch(selected)
	}
	return nil
}
