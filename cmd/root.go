package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"github.com/smallstep/cli-utils/step"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "step-tpm-plugin",
	Short: "🔑 `step` plugin for interacting with TPMs.",
	Long:  `🔑 step plugin for interacting with TPMs.`,
}

// Execute adds all child commands to the root command and prepares the step
// environment. It then invokes the root command in style with fangs.
func Execute() {
	// add all child commands
	rootCmd.AddCommand(
		NewInfoCommand(),
		NewEKCommand(),
		NewAKCommand(),
		NewKeyCommand(),
		NewRandomCommand(),
		NewSimulatorCommand(),
		NewVersionCommand(),
	)

	// ensure step environment is prepared before every command
	rootCmd.PersistentPreRunE = func(*cobra.Command, []string) error {
		if err := step.Init(); err != nil {
			return fmt.Errorf("failed initializing step environment: %w", err)
		}
		return nil
	}

	// execute the command
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}
