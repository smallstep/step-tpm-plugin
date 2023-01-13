package keys

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm"
)

func NewCreateKeyCommand() *cobra.Command {
	const (
		long = `create a TPM key. Specify --ak=akName to attest the key.
`
		short = "create a TPM key"
	)

	cmd := command.New("create", short, long, runCreateKey, []command.Preparer{command.RequireTPMWithStorage}, nil)

	cmd.Args = cobra.RangeArgs(0, 1)

	flag.Add(cmd,
		flag.StorageFile(),
		flag.StorageDirectory(),
		flag.JSON(),
		flag.Device(),
		flag.String{
			Name:        "ak",
			Description: "Name of the AK to attest new key with",
		},
	)

	return cmd
}

func runCreateKey(ctx context.Context) error {

	var (
		t      = tpm.FromContext(ctx)
		name   = flag.FirstArg(ctx)
		akName = flag.GetString(ctx, "ak")
	)

	var (
		key tpm.Key
		err error
	)

	if akName == "" {
		// create a key without attesting to an AK
		if key, err = t.CreateKey(ctx, name); err != nil {
			return fmt.Errorf("creating key failed: %w", err)
		}
	} else {
		// create a key attested by AK
		if key, err = t.AttestKey(ctx, akName, name); err != nil {
			return fmt.Errorf("creating attested key failed: %w", err)
		}
	}

	t1 := table.NewWriter()
	t1.SetOutputMirror(os.Stdout)
	t1.AppendHeader(table.Row{"Name", "Data"})
	t1.AppendRow(table.Row{key.Name, len(key.Data)})
	t1.Render()

	return nil
}
