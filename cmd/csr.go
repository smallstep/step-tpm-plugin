package cmd

import (
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/cmd/csr"
	"github.com/smallstep/step-tpm-plugin/internal/command"
)

func NewCSRCommand() *cobra.Command {
	const (
		long = `subcommand for CSRs.
`
		short = "subcommand for CSRs"
	)

	cmd := command.New("csr <command>", short, long, nil, nil, nil)

	cmd.AddCommand(
		csr.NewSignCSRCommand(),
	)

	return cmd
}
