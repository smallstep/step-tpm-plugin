package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/smallstep/cli-utils/step"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "step-tpm-plugin",
	Short: "🔑 `step` plugin for interacting with TPMs. ",
	Long: `🔑 step plugin for interacting with TPMs. 
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// initialize step environment.
	if err := step.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.AddCommand(
		NewInfoCommand(),
		NewEKCommand(),
		NewAKCommand(),
		NewKeyCommand(),
		NewRandomCommand(),
		NewSimulatorCommand(),
		NewVersionCommand(),
	)
}
