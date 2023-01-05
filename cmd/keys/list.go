package keys

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"

	"github.com/smallstep/step-tpm-plugin/pkg/tpm"
)

// TODO(hs): move to internal?

func NewListKeysCommand() *cobra.Command {
	const (
		long = `loooooooong list description.
`
		short = "short list description"
	)

	cmd := command.New("list", short, long, runListKeys, []command.Preparer{command.RequireTPMStore, command.RequireTPM}, nil)

	cmd.Args = cobra.NoArgs

	flag.Add(cmd,
		flag.JSON(),
		flag.Device(),
		flag.StorageFile(),
	)

	return cmd
}

func runListKeys(ctx context.Context) error {

	var (
		t    = tpm.FromContext(ctx)
		json = flag.GetBool(ctx, flag.FlagJSON)
	)

	_ = json

	keys, err := t.ListKeys(ctx)
	if err != nil {
		return err
	}

	for _, key := range keys {
		fmt.Println(key.Name, key.Data)
	}

	return nil
}
