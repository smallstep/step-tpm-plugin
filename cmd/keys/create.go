package keys

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"go.step.sm/crypto/tpm"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/internal/render"
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
		flag.AK(),
		flag.Int{
			Name:        "size",
			Description: "Size of key to create",
			Default:     2048,
		},
		flag.String{
			Name:        "kty",
			Description: "Type of key to create",
			Default:     "RSA",
		},
	)

	return cmd
}

func runCreateKey(ctx context.Context) error {
	var (
		t      = tpm.FromContext(ctx)
		name   = flag.FirstArg(ctx)
		json   = flag.GetBool(ctx, flag.FlagJSON)
		akName = flag.GetString(ctx, flag.FlagAK)
		size   = flag.GetInt(ctx, "size")
		kty    = flag.GetString(ctx, "kty")
	)

	// TODO: validate size, combined with (valid) key algorithms

	var (
		key *tpm.Key
		err error
	)

	if akName == "" {
		// create a key without attesting to an AK
		config := tpm.CreateKeyConfig{
			Algorithm: kty,
			Size:      size,
		}
		if key, err = t.CreateKey(ctx, name, config); err != nil {
			return fmt.Errorf("creating key failed: %w", err)
		}
	} else {
		// create a key attested by AK
		config := tpm.AttestKeyConfig{
			Algorithm: kty,
			Size:      size,
		}
		if key, err = t.AttestKey(ctx, akName, name, config); err != nil {
			return fmt.Errorf("creating attested key failed: %w", err)
		}
	}

	if json {
		return render.JSON(os.Stdout, key)
	}

	t1 := table.NewWriter()
	t1.SetOutputMirror(os.Stdout)
	t1.AppendHeader(table.Row{"Name", "Data"})
	t1.AppendRow(table.Row{key.Name(), len(key.Data())})
	t1.Render()

	return nil
}
