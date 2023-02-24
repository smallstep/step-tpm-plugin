package keys

import (
	"context"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"go.step.sm/crypto/tpm"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/internal/render"
)

func NewListKeysCommand() *cobra.Command {
	const (
		long = `list TPM key details.
`
		short = "list TPM key details"
	)

	cmd := command.New("list", short, long, runListKeys, []command.Preparer{command.RequireTPMWithStorage}, nil)

	cmd.Args = cobra.NoArgs

	flag.Add(cmd,
		flag.StorageFile(),
		flag.StorageDirectory(),
		flag.JSON(),
		flag.Device(),
		flag.AK(),
	)

	return cmd
}

func runListKeys(ctx context.Context) error {
	var (
		t    = tpm.FromContext(ctx)
		json = flag.GetBool(ctx, flag.FlagJSON)
		ak   = flag.GetString(ctx, flag.FlagAK)
	)

	var keys []*tpm.Key
	var err error
	if ak != "" {
		keys, err = t.GetKeysAttestedBy(ctx, ak)
	} else {
		keys, err = t.ListKeys(ctx)
	}
	if err != nil {
		return err
	}

	if json {
		return render.JSON(os.Stdout, keys)
	}

	t1 := table.NewWriter()
	t1.SetOutputMirror(os.Stdout)
	t1.AppendHeader(table.Row{"Name", "Data length"})
	for _, key := range keys {
		t1.AppendRow(table.Row{key.Name(), len(key.Data())})
	}
	t1.Render()

	return nil
}
