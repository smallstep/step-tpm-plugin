package ak

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/internal/render"
	"go.step.sm/crypto/tpm"
)

func NewCreateAKCommand() *cobra.Command {
	const (
		long = `create an AK.
`
		short = "create an AK"
	)

	cmd := command.New("create", short, long, runCreateAK, []command.Preparer{command.RequireTPMWithStorage}, nil)

	cmd.Args = cobra.RangeArgs(0, 1)

	flag.Add(cmd,
		flag.StorageFile(),
		flag.StorageDirectory(),
		flag.Device(),
		flag.JSON(),
	)

	return cmd
}

func runCreateAK(ctx context.Context) error {
	var (
		t    = tpm.FromContext(ctx)
		name = flag.FirstArg(ctx)
		json = flag.GetBool(ctx, flag.FlagJSON)
	)

	ak, err := t.CreateAK(ctx, name)
	if err != nil {
		return fmt.Errorf("creating AK failed: %w", err)
	}

	if json {
		return render.JSON(os.Stdout, ak)
	}

	t1 := table.NewWriter()
	t1.SetOutputMirror(os.Stdout)
	t1.AppendHeader(table.Row{"Name", "Data"})
	t1.AppendRow(table.Row{ak.Name(), len(ak.Data())})
	t1.Render()

	return nil
}
