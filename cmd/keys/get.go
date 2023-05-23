package keys

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"go.step.sm/crypto/pemutil"
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
		flag.PEM(),
		flag.Bundle(),
		flag.Blob(),
		flag.Private(),
		flag.Public(),
	)

	return cmd
}

func runGetKey(ctx context.Context) error {
	var (
		t             = tpm.FromContext(ctx)
		json          = flag.GetBool(ctx, flag.FlagJSON)
		name          = flag.FirstArg(ctx)
		outputPEM     = flag.GetBool(ctx, flag.FlagPEM)
		bundle        = flag.GetBool(ctx, flag.FlagBundle)
		outputBlob    = flag.GetBool(ctx, flag.FlagBlob)
		outputPrivate = flag.GetBool(ctx, flag.FlagPrivate)
		outputPublic  = flag.GetBool(ctx, flag.FlagPublic)
	)

	key, err := t.GetKey(ctx, name)
	if err != nil {
		return fmt.Errorf("getting key failed: %w", err)
	}

	// TODO(hs): add option to write to file?
	// TODO(hs): add flag to output hex?
	if outputPEM {
		chain := key.CertificateChain()
		if len(chain) == 0 {
			fmt.Println("no certificate available")
			return nil
		}
		if len(chain) > 1 && bundle {
			var b []byte
			for _, crt := range chain {
				b = append(b, pem.EncodeToMemory(&pem.Block{
					Type:  "CERTIFICATE",
					Bytes: crt.Raw,
				})...)
			}
			_, err = fmt.Println(string(b))
			return err
		}
		b, err := pemutil.Serialize(chain[0])
		if err != nil {
			return fmt.Errorf("failed serializing certificate: %w", err)
		}
		_, err = fmt.Println(string(pem.EncodeToMemory(b)))
		return err
	}

	if outputBlob {
		blobs, err := key.Blobs(ctx)
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
		return render.JSON(os.Stdout, key)
	}

	// TODO(hs): dumping the raw data isn't the most useful thing to do
	fmt.Println(base64.StdEncoding.EncodeToString(key.Data()))

	return nil
}
