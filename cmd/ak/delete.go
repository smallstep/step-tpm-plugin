package ak

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm"
)

func NewDeleteAKCommand() *cobra.Command {
	const (
		long = `delete an AK.
`
		short = "delete an AK"
	)

	cmd := command.New("delete", short, long, runDeleteAK, []command.Preparer{command.RequireTPM}, nil)

	cmd.Args = cobra.ExactArgs(1)

	flag.Add(cmd,
		flag.StorageFile(),
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
