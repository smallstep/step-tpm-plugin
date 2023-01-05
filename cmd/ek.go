package cmd

import (
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/cmd/ek"
	"github.com/smallstep/step-tpm-plugin/internal/command"
)

func NewEKCommand() *cobra.Command {
	const (
		long = `subcommand for EKs.
`
		short = "subcommand for EKs"
	)

	cmd := command.New("ek <command>", short, long, nil, nil, nil) // TODO: actually do something

	cmd.AddCommand(
		ek.NewGetEKCommand(),
	)

	return cmd
}
