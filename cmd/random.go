package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm"
)

func NewRandomCommand() *cobra.Command {
	const (
		long = `subcommand for generating random data.
`
		short = "subcommand for for generating random data"
	)

	cmd := command.New("random <command>", short, long, runRandom, []command.Preparer{command.RequireTPM}, nil)

	flag.Add(cmd,
		flag.Device(), // TOOD(hs): currently unused here. Should affect (lazy) instantiation of the TPM.
		flag.StorageFile(),
	)

	return cmd
}

func runRandom(ctx context.Context) error {

	var (
		t = tpm.FromContext(ctx)
	)

	r, err := t.GenerateRandom(ctx, 32)
	if err != nil {
		return fmt.Errorf("failed generating random: %w", err)
	}

	fmt.Println(r)

	return nil
}
