package cmd

import (
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/cmd/keys"
	"github.com/smallstep/step-tpm-plugin/internal/command"
)

func NewKeyCommand() *cobra.Command {
	const (
		long = `subcommand for managing TPM keys.
`
		short = "subcommand for managing TPM keys"
	)

	cmd := command.New("key <command>", short, long, nil, nil, nil)

	cmd.AddCommand(
		keys.NewCreateKeyCommand(),
		keys.NewListKeysCommand(),
		keys.NewGetKeyCommand(),
		keys.NewDeleteKeyCommand(),
	)

	// TODO: look into making keys persist (using EvictControl) in the TPM? How to go about deletion in that case?

	return cmd
}
