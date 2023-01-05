package cmd

import (
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/cmd/ak"
	"github.com/smallstep/step-tpm-plugin/internal/command"
)

func NewAKCommand() *cobra.Command {
	const (
		long = `loooooooong ak command description.
`
		short = "short ak command description"
	)

	cmd := command.New("ak <command>", short, long, nil, nil, nil)

	cmd.AddCommand(
		ak.NewCreateAKCommand(),
		ak.NewListAKCommand(),
		ak.NewGetAKCommand(),
		ak.NewDeleteAKCommand(),
	)

	// TODO: add AK certificate command (and handle storage for that too?), ...
	// TODO: look into making keys persistent in the TPM? How to go about deletion in that case?

	return cmd
}
