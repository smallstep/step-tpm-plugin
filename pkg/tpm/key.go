package tpm

import (
	"context"
	"crypto"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/google/go-attestation/attest"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm/internal/key"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm/storage"
)

func (t *TPM) CreateKey(ctx context.Context, name string) (*storage.Key, error) {

	if err := t.Open(ctx); err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	if name == "" {
		nameHex := make([]byte, 5)
		if n, err := rand.Read(nameHex); err != nil || n != len(nameHex) {
			return nil, fmt.Errorf("rand.Read() failed with %d/%d bytes read and error: %v", n, len(nameHex), err)
		}
		name = fmt.Sprintf("%x", nameHex)
	}

	prefixedKeyName := fmt.Sprintf("app-%s", name)

	data, err := key.Create(t.deviceName, prefixedKeyName) // TODO: additional parameters
	if err != nil {
		return nil, fmt.Errorf("failed creating key: %w", err)
	}

	storedKey := &storage.Key{
		Name: name,
		Data: data,
	}

	if err := t.store.AddKey(storedKey); err != nil {
		return nil, err
	}

	if err := t.store.Persist(); err != nil {
		return nil, fmt.Errorf("error persisting to storage: %w", err)
	}

	return storedKey, nil
}

// TODO: every interaction with the actual TPM now opens the "connection" when required, then
// closes it when the operation is done. Can we reuse one open "connection" to the TPM for
// multiple operations reliably? What makes it harder is that now all operations are implemented
// by go-attestation, so it might come down to replicating a lot of that logic. It could involve
// checking multiple locks and/or pointers and instantiating when required. Opening and closing
// on-demand is the simplest way and safe to do for now, though.
func (t *TPM) AttestKey(ctx context.Context, akName, name string) (*storage.Key, error) {

	if err := t.Open(ctx); err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	at, err := attest.OpenTPM(t.attestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer at.Close()

	ak, err := t.store.GetAK(akName)
	if err != nil {
		return nil, err
	}

	loadedAK, err := at.LoadAK(ak.Data)
	if err != nil {
		return nil, err
	}
	defer loadedAK.Close(at)

	if name == "" {
		nameHex := make([]byte, 5)
		if n, err := rand.Read(nameHex); err != nil || n != len(nameHex) {
			return nil, fmt.Errorf("rand.Read() failed with %d/%d bytes read and error: %v", n, len(nameHex), err)
		}
		name = fmt.Sprintf("%x", nameHex)
	}

	prefixedKeyName := fmt.Sprintf("app-%s", name)

	keyConfig := &attest.KeyConfig{
		// TODO: provide values (through flags) for algorithm, size, name, prefix, qualifying data?
		Algorithm:      attest.RSA,
		Size:           2048,
		QualifyingData: nil, // TODO: needs value for ACME `device-attest-01`
		Name:           prefixedKeyName,
	}

	key, err := at.NewKey(loadedAK, keyConfig)
	if err != nil {
		return nil, err
	}
	defer key.Close()

	data, err := key.Marshal()
	if err != nil {
		return nil, err
	}

	storedKey := &storage.Key{
		Name: name,
		Data: data,
	}

	if err := t.store.AddKey(storedKey); err != nil {
		return nil, err
	}

	if err := t.store.Persist(); err != nil {
		return nil, fmt.Errorf("error persisting to storage: %w", err)
	}

	return storedKey, nil
}

func (t *TPM) GetKey(ctx context.Context, name string) (*storage.Key, error) {

	if err := t.Open(ctx); err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	return t.store.GetKey(name)
}

func (t *TPM) ListKeys(ctx context.Context) ([]*storage.Key, error) {

	if err := t.Open(ctx); err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	return t.store.ListKeys()
}

func (t *TPM) GetKeyCertificationParameters(ctx context.Context, name string) (attest.CertificationParameters, error) {

	result := attest.CertificationParameters{}
	if err := t.Open(ctx); err != nil {
		return result, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	at, err := attest.OpenTPM(t.attestConfig)
	if err != nil {
		return result, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer at.Close()

	key, err := t.store.GetKey(name)
	if err != nil {
		return result, err
	}

	loadedKey, err := at.LoadKey(key.Data) // TODO: store the attestation parameters in the keystore instead too? That makes this operation simpler
	if err != nil {
		return attest.CertificationParameters{}, err
	}

	return loadedKey.CertificationParameters(), nil
}

func (t *TPM) DeleteKey(ctx context.Context, name string) error {
	if err := t.Open(ctx); err != nil {
		return fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	at, err := attest.OpenTPM(t.attestConfig)
	if err != nil {
		return fmt.Errorf("failed opening TPM: %w", err)
	}
	defer at.Close()

	key, err := t.store.GetKey(name)
	if err != nil {
		return fmt.Errorf("failed loading key: %w", err)
	}

	if err := at.DeleteKey(key.Data); err != nil {
		return fmt.Errorf("failed deleting key: %w", err)
	}

	if err := t.store.DeleteKey(name); err != nil {
		return fmt.Errorf("error deleting key: %w", err)
	}

	if err := t.store.Persist(); err != nil {
		return fmt.Errorf("error persisting storage: %w", err)
	}

	return nil
}

// signer implements crypto.Signer backed by a TPM key
type signer struct {
	tpm    *TPM
	key    *storage.Key
	public crypto.PublicKey
}

func (s *signer) Public() crypto.PublicKey {
	return s.public
}

func (s *signer) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {

	ctx := context.Background()
	if err := s.tpm.Open(ctx); err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer s.tpm.Close(ctx, false)

	at, err := attest.OpenTPM(s.tpm.attestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer at.Close()

	loadedKey, err := at.LoadKey(s.key.Data)
	if err != nil {
		return nil, err
	}
	defer loadedKey.Close()

	priv, err := loadedKey.Private(s.public)
	if err != nil {
		return nil, err
	}

	var signer crypto.Signer
	var ok bool
	if signer, ok = priv.(crypto.Signer); !ok {
		return nil, errors.New("error getting TPM private key as crypto.Signer")
	}

	return signer.Sign(rand, digest, opts)
}

// GetSigner ...
//
// TODO: conclude whether AKs should also be usable as signers?
func (t *TPM) GetSigner(ctx context.Context, name string) (crypto.Signer, error) {

	if err := t.Open(ctx); err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	at, err := attest.OpenTPM(t.attestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer at.Close()

	key, err := t.store.GetKey(name)
	if err != nil {
		return nil, err
	}

	loadedKey, err := at.LoadKey(key.Data)
	if err != nil {
		return nil, err
	}
	defer loadedKey.Close()

	return &signer{
		tpm:    t,
		key:    key,
		public: loadedKey.Public(),
	}, nil
}
