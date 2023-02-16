package ak

import (
	"context"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/internal/render"
	"go.step.sm/crypto/tpm"
)

// TODO(hs): move to internal?

func NewListAKCommand() *cobra.Command {
	const (
		long = `list AK details.
`
		short = "list AK details"
	)

	cmd := command.New("list", short, long, runListAK, []command.Preparer{command.RequireTPMWithStorage}, nil)

	flag.Add(cmd,
		flag.StorageFile(),
		flag.StorageDirectory(),
		flag.Device(),
		flag.JSON(),
	)

	return cmd
}

func runListAK(ctx context.Context) error {
	var (
		t    = tpm.FromContext(ctx)
		json = flag.GetBool(ctx, flag.FlagJSON)
	)

	aks, err := t.ListAKs(ctx)
	if err != nil {
		return err
	}

	if json {
		return render.JSON(os.Stdout, aks)
	}

	t1 := table.NewWriter()
	t1.SetOutputMirror(os.Stdout)
	t1.AppendHeader(table.Row{"Name", "Data length"})
	for _, ak := range aks {
		t1.AppendRow(table.Row{ak.Name(), len(ak.Data())})
	}
	t1.Render()

	return nil
}
