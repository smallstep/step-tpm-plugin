package keys

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/pkg/tpm"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
)

func NewDeleteKeyCommand() *cobra.Command {
	const (
		long = `loooooooong delete description.
`
		short = "short delete description"
	)

	cmd := command.New("delete", short, long, runDeleteKey, []command.Preparer{command.RequireTPMStore, command.RequireTPM}, nil)

	cmd.Args = cobra.RangeArgs(0, 1)

	flag.Add(cmd,
		flag.JSON(),
		flag.Device(),
		flag.StorageFile(),
	)

	return cmd
}

func runDeleteKey(ctx context.Context) error {

	var (
		t    = tpm.FromContext(ctx)
		name = flag.FirstArg(ctx)
	)

	if err := t.DeleteKey(ctx, name); err != nil {
		return fmt.Errorf("error deleting key: %w", err)
	}

	return nil
}
