package command

import (
	"context"
	"fmt"

	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"go.step.sm/crypto/tpm"
	"go.step.sm/crypto/tpm/storage"
)

// fallbackTPMStore ensures a NOOP TPMStore is available in the context.
func fallbackTPMStore(ctx context.Context) (context.Context, error) {
	ctx = storage.BlackHoleContext(ctx) // set default storage in context
	return ctx, nil
}

func RequireTPMWithoutStorage(ctx context.Context) (context.Context, error) {
	var (
		deviceName = flag.GetString(ctx, flag.FlagDeviceName) // TODO(hs): it feels a bit messy to rely on the flag here; can we improve?
	)

	pt, err := tpm.New(tpm.WithDeviceName(deviceName))
	if err != nil {
		return nil, fmt.Errorf("failed creating TPM: %w", err)
	}

	ctx = tpm.NewContext(ctx, pt)

	return ctx, nil
}

func RequireTPMWithStorage(ctx context.Context) (context.Context, error) {
	var (
		deviceName       = flag.GetString(ctx, flag.FlagDeviceName) // TODO(hs): it feels a bit messy to rely on the flag here; can we improve?
		storageFile      = flag.GetString(ctx, flag.FlagStorageFile)
		storageDirectory = flag.GetString(ctx, flag.FlagStorageDirectory)
	)

	var store storage.TPMStore
	switch {
	case storageFile != flag.StorageFile().Default: // TODO: not entirely happy with the logic, but it works
		store = storage.NewFilestore(storageFile)
	default:
		store = storage.NewDirstore(storageDirectory)
	}

	pt, err := tpm.New(tpm.WithDeviceName(deviceName), tpm.WithStore(store))
	if err != nil {
		return nil, fmt.Errorf("failed creating TPM: %w", err)
	}

	ctx = tpm.NewContext(ctx, pt)

	return ctx, nil
}
