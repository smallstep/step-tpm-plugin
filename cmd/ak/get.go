package ak

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"go.step.sm/crypto/tpm"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/internal/render"
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
		flag.Blob(),
		flag.Private(),
		flag.Public(),
	)

	return cmd
}

func runGetAK(ctx context.Context) error {
	var (
		t             = tpm.FromContext(ctx)
		json          = flag.GetBool(ctx, flag.FlagJSON)
		name          = flag.FirstArg(ctx)
		outputBlob    = flag.GetBool(ctx, flag.FlagBlob)
		outputPrivate = flag.GetBool(ctx, flag.FlagPrivate)
		outputPublic  = flag.GetBool(ctx, flag.FlagPublic)
	)

	ak, err := t.GetAK(ctx, name)
	if err != nil {
		return fmt.Errorf("getting AK failed: %w", err)
	}

	// TODO(hs): add option to write to file?
	// TODO(hs): add flag to output hex?

	if outputBlob {
		blobs, err := ak.Blobs(ctx)
		if err != nil {
			return fmt.Errorf("failed getting key blobs: %w", err)
		}

		switch {
		case outputPrivate:
			private, err := blobs.Private()
			if err != nil {
				return fmt.Errorf("failed getting private: %w", err)
			}
			fmt.Println(string(private))
		case outputPublic:
			public, err := blobs.Public()
			if err != nil {
				return fmt.Errorf("failed getting public: %w", err)
			}
			fmt.Println(string(public))
		default:
			return errors.New("--private or --public required")
		}

		return nil
	}

	if json {
		return render.JSON(os.Stdout, ak)
	}

	// TODO(hs): dumping the raw data isn't the most useful thing to do
	fmt.Println(base64.StdEncoding.EncodeToString(ak.Data()))

	return nil
}
