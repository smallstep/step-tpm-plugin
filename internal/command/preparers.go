package command

import (
	"context"
	"fmt"

	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm/storage"
)

// fallbackTPMStore ensures a NOOP TPMStore is available in the context.
func fallbackTPMStore(ctx context.Context) (context.Context, error) {
	ctx = storage.BlackHoleContext(ctx) // set default storage in context
	return ctx, nil
}

func RequireTPM(ctx context.Context) (context.Context, error) {

	var (
		deviceName  = flag.GetString(ctx, flag.FlagDeviceName) // TODO(hs): it feels a bit messy to rely on the flag here; can we improve?
		storageFile = flag.GetString(ctx, flag.FlagStorageFile)
		filestore   = storage.NewFilestore(storageFile) // TODO: is the file sufficient? It should allow absolute paths already.
	)

	if err := filestore.Load(); err != nil { // TODO: can fail if the file doesn't exist; create on first usage?
		return nil, fmt.Errorf("failed loading filestore: %w", err)
	}

	pt, err := tpm.New(tpm.WithDeviceName(deviceName), tpm.WithStore(filestore))
	if err != nil {
		return nil, fmt.Errorf("failed creating TPM: %w", err)
	}

	ctx = tpm.NewContext(ctx, pt)

	return ctx, nil
}
