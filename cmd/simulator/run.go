package simulator

import (
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
)

func NewRunCommand() *cobra.Command {
	const (
		long  = `Run the TPM simulator.`
		short = "Run the TPM simulator."
	)

	cmd := command.New("run", short, long, runSimulator, nil, nil)

	flag.Add(cmd,
		flag.Socket(),
		flag.Seed(),
		flag.Verbose(),
	)

	return cmd
}
