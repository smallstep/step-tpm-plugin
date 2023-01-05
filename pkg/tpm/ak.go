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
			return nil, err
		}
		name = fmt.Sprintf("%x", nameHex)
	}

	prefixedName := fmt.Sprintf("ak-%s", name)

	akConfig := attest.AKConfig{
		//Name: "test-this", // NOTE: specifying the name and trying to create this multiple times will error
		Name: prefixedName,
	}
	ak, err := at.NewAK(&akConfig)
	if err != nil {
		return nil, err
	}
	defer ak.Close(at)

	fmt.Println(ak)

	b, err := ak.Marshal()
	if err != nil {
		return nil, err
	}

	fmt.Println(b)

	fmt.Println(ak.AttestationParameters())

	storedAK := &storage.AK{
		// TODO: name will only work on Windows; on Linux there's no way to specify this,
		// there'll only be a blob data. Provide our own identifier instead? Or include a hash
		// as part of the name, so that it can be looked up? Hash of the public key? We want to
		// be able to easily specify certain names/files for the keys so that they can be used
		// for attesting other keys, then new keys can be created and attested, and the keys can
		// then be more easily used to be loaded for signing operations.
		Name: name,
		Data: b,
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
