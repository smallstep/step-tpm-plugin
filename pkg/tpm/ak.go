package tpm

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/google/go-attestation/attest"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm/storage"
)

func (t *TPM) CreateAK(ctx context.Context, name string) (*storage.AK, error) { // TODO: return information about AK too?

	if err := t.Open(ctx); err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	at, err := attest.OpenTPM(t.attestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer at.Close()

	if name == "" {
		// TODO: decouple the TPM key name from the name recorded in the storage? This might
		// make it easier to work with the key names as a user; the TPM key name would be abstracted
		// away. The key name in the storage can be different from the key stored with the key (which,
		// to be far, isn't even used on Linux TPMs)
		nameHex := make([]byte, 5)
		if n, err := rand.Read(nameHex); err != nil || n != len(nameHex) {
			return nil, fmt.Errorf("rand.Read() failed with %d/%d bytes read and error: %v", n, len(nameHex), err)
		}
		name = fmt.Sprintf("%x", nameHex)
	}

	prefixedName := fmt.Sprintf("ak-%s", name)

	akConfig := attest.AKConfig{
		Name: prefixedName,
	}
	ak, err := at.NewAK(&akConfig)
	if err != nil {
		return nil, err
	}
	defer ak.Close(at)

	fmt.Println(ak)

	data, err := ak.Marshal()
	if err != nil {
		return nil, err
	}

	fmt.Println(ak.AttestationParameters())

	storedAK := &storage.AK{
		Name: name,
		Data: data,
	}

	if err := t.store.AddAK(storedAK); err != nil {
		return nil, err
	}

	if err := t.store.Persist(); err != nil {
		return nil, err
	}

	return storedAK, nil
}

func (t *TPM) GetAK(ctx context.Context, name string) (*storage.AK, error) {

	if err := t.Open(ctx); err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	return t.store.GetAK(name)
}

func (t *TPM) ListAKs(ctx context.Context) ([]*storage.AK, error) {

	if err := t.Open(ctx); err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	return t.store.ListAKs()
}

func (t *TPM) DeleteAK(ctx context.Context, name string) error {

	if err := t.Open(ctx); err != nil {
		return fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	if err := t.store.DeleteAK(name); err != nil {
		return fmt.Errorf("error deleting AK: %w", err)
	}

	if err := t.store.Persist(); err != nil {
		return fmt.Errorf("error persisting storage: %w", err)
	}

	return nil
}
