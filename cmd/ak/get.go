package ak

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/pkg/tpm"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
)

func NewGetAKCommand() *cobra.Command {
	const (
		long = `loooooooong get description.
`
		short = "short get description"
	)

	cmd := command.New("get", short, long, runGetAK, []command.Preparer{command.RequireTPM}, nil)

	cmd.Args = cobra.ExactArgs(1)

	flag.Add(cmd,
		flag.JSON(),
		flag.Device(),
		flag.StorageFile(),
	)

	return cmd
}

func runGetAK(ctx context.Context) error {

	var (
		t    = tpm.FromContext(ctx)
		json = flag.GetBool(ctx, flag.FlagJSON)
		name = flag.FirstArg(ctx)
	)

	_ = json

	key, err := t.GetAK(ctx, name)
	if err != nil {
		return fmt.Errorf("getting AK failed: %w", err)
	}

	fmt.Println(key)

	return nil
}
