package keys

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/pkg/tpm"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
)

func NewGetKeyCommand() *cobra.Command {
	const (
		long = `print TPM key base64 blob data.
`
		short = "print TPM key base64 blob data"
	)

	cmd := command.New("get", short, long, runGetKey, []command.Preparer{command.RequireTPM}, nil) // TODO: rename this blob?

	cmd.Args = cobra.ExactArgs(1)

	flag.Add(cmd,
		flag.JSON(),
		flag.Device(),
		flag.StorageFile(),
	)

	return cmd
}

func runGetKey(ctx context.Context) error {

	var (
		t    = tpm.FromContext(ctx)
		json = flag.GetBool(ctx, flag.FlagJSON)
		name = flag.FirstArg(ctx)
	)

	_ = json

	key, err := t.GetKey(ctx, name)
	if err != nil {
		return fmt.Errorf("getting key failed: %w", err)
	}

	// t1 := table.NewWriter()
	// t1.SetOutputMirror(os.Stdout)
	// t1.AppendHeader(table.Row{"Name", "Data"})
	// t1.AppendRow(table.Row{key.Name, len(key.Data)})
	// t1.Render()

	fmt.Println(base64.StdEncoding.EncodeToString(key.Data))

	return nil
}
