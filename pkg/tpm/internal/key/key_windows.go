//go:build windows
// +build windows

package key

import (
	"fmt"

	"github.com/google/go-attestation/attest"
)

func create(_, keyName string) ([]byte, error) {

	pcp, err := openPCP()
	if err != nil {
		return nil, fmt.Errorf("failed to open PCP: %w", err)
	}
	defer pcp.Close()

	config := &KeyConfig{Algorithm: RSA, Size: 2048} // TODO: additional parameters

	hnd, pub, blob, err := pcp.NewKey(keyName, config)
	if err != nil {
		return nil, fmt.Errorf("pcp failed to mint application key: %v", err)
	}

	_, _ = hnd, blob

	// tpmPub, err := tpm2.DecodePublic(pub)
	// if err != nil {
	// 	return nil, fmt.Errorf("decode public key: %v", err)
	// }

	// pubKey, err := tpmPub.Key()
	// if err != nil {
	// 	return nil, fmt.Errorf("access public key: %v", err)
	// }

	out := serializedKey{
		Encoding:   keyEncodingOSManaged,
		TPMVersion: attest.TPMVersion20,
		Name:       keyName,
		Public:     pub,
	}

	return out.Serialize()
}
