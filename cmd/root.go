package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "step-tpm-plugin",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },

	// TODO(hs): make this more composable, so that multiple operations can be
	// performed in preparation of (all) commands?
	// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

	// 	// open connection to TPM and add it to the context for commands to use
	// 	t, err := tpm.Open()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	cmd.SetContext(tpm.NewContext(cmd.Context(), t))

	// 	return nil
	// },

	// TODO(hs): make this more composable, so that multiple operations can be
	// performed to clean up after running (all) commands?
	// PersistentPostRunE: func(cmd *cobra.Command, args []string) error {

	// 	// get the TPM and properly close it after running
	// 	t := tpm.FromContext(cmd.Context())
	// 	if t == nil {
	// 		return nil
	// 	}

	// 	return t.Close()
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.step-tpm-plugin.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(
		NewInfoCommand(),
		NewEKCommand(),
		NewAKCommand(),
		NewKeyCommand(),
		NewCSRCommand(),
	)

}
