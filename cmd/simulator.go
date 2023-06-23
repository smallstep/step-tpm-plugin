package cmd

import (
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/cmd/simulator"
	"github.com/smallstep/step-tpm-plugin/internal/command"
)

func NewSimulatorCommand() *cobra.Command {
	const (
		long  = `subcommand for TPM simulator.`
		short = "subcommand for TPM simulator"
	)

	cmd := command.New("simulator <command>", short, long, nil, nil, nil)

	cmd.AddCommand(
		simulator.NewRunCommand(),
	)

	return cmd
}
