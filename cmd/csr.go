package cmd

import (
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/cmd/csr"
	"github.com/smallstep/step-tpm-plugin/internal/command"
)

func NewCSRCommand() *cobra.Command {
	const (
		long = `loooooooong csr command description.
`
		short = "short csr command description"
	)

	cmd := command.New("csr <command>", short, long, nil, nil, nil)

	cmd.AddCommand(
		csr.NewSignCSRCommand(),
	)

	return cmd
}
