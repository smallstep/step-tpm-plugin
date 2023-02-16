package keys

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"go.step.sm/crypto/tpm"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/internal/render"
)

func NewGetKeyCommand() *cobra.Command {
	const (
		long = `print TPM key base64 blob data.
`
		short = "print TPM key base64 blob data"
	)

	cmd := command.New("get", short, long, runGetKey, []command.Preparer{command.RequireTPMWithStorage}, nil) // TODO: rename this blob?

	cmd.Args = cobra.ExactArgs(1)

	flag.Add(cmd,
		flag.StorageFile(),
		flag.StorageDirectory(),
		flag.JSON(),
		flag.Device(),
	)

	return cmd
}

func runGetKey(ctx context.Context) error {
	var (
		t    = tpm.FromContext(ctx)
		json = flag.GetBool(ctx, flag.FlagJSON)
		name = flag.FirstArg(ctx)
	)

	key, err := t.GetKey(ctx, name)
	if err != nil {
		return fmt.Errorf("getting key failed: %w", err)
	}

	if json {
		return render.JSON(os.Stdout, key)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(key.Data()))

	return nil
}
