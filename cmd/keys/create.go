package keys

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm/storage"
)

func NewCreateKeyCommand() *cobra.Command {
	const (
		long = `loooooooong create description.
`
		short = "short create description"
	)

	cmd := command.New("create", short, long, runCreateKey, []command.Preparer{command.RequireTPMStore, command.RequireTPM}, nil)

	cmd.Args = cobra.RangeArgs(0, 1)

	flag.Add(cmd,
		flag.JSON(),
		flag.Device(),
		flag.StorageFile(),
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
		key *storage.Key
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

	fmt.Println(key)

	return nil
}
