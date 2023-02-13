package ak

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"go.step.sm/crypto/tpm"
)

func NewDeleteAKCommand() *cobra.Command {
	const (
		long = `delete an AK.
`
		short = "delete an AK"
	)

	cmd := command.New("delete", short, long, runDeleteAK, []command.Preparer{command.RequireTPMWithStorage}, nil)

	cmd.Args = cobra.ExactArgs(1)

	flag.Add(cmd,
		flag.StorageFile(),
		flag.StorageDirectory(),
		flag.Device(),
		flag.JSON(),
	)

	return cmd
}

func runDeleteAK(ctx context.Context) error {
	var (
		t    = tpm.FromContext(ctx)
		name = flag.FirstArg(ctx)
	)

	if err := t.DeleteAK(ctx, name); err != nil {
		return fmt.Errorf("deleting AK failed: %w", err)
	}

	return nil
}
