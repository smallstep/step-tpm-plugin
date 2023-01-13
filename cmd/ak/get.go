package ak

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/pkg/tpm"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
)

func NewGetAKCommand() *cobra.Command {
	const (
		long = `print AK base64 blob data.
`
		short = "print AK base64 blob data"
	)

	cmd := command.New("get", short, long, runGetAK, []command.Preparer{command.RequireTPMWithStorage}, nil)

	cmd.Args = cobra.ExactArgs(1)

	flag.Add(cmd,
		flag.StorageFile(),
		flag.StorageDirectory(),
		flag.JSON(),
		flag.Device(),
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

	ak, err := t.GetAK(ctx, name)
	if err != nil {
		return fmt.Errorf("getting AK failed: %w", err)
	}

	// t1 := table.NewWriter()
	// t1.SetOutputMirror(os.Stdout)
	// t1.AppendHeader(table.Row{"Name", "Data"})
	// t1.AppendRow(table.Row{ak.Name, len(ak.Data)})
	// t1.Render()

	fmt.Println(base64.StdEncoding.EncodeToString(ak.Data))

	return nil
}
