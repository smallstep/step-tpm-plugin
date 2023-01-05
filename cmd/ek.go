package cmd

import (
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
)

func NewEKCommand() *cobra.Command {
	const (
		long = `loooooooong ek command description.
`
		short = "short ek command description"
	)

	cmd := command.New("ek <command>", short, long, nil, nil, nil) // TODO: actually do something

	return cmd
}
