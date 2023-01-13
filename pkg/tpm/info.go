package tpm

import (
	"context"
	"fmt"

	"github.com/google/go-attestation/attest"
)

type Info struct {
	Version              uint8 // TODO: alias the `attest` types instead?
	Interface            uint8
	VendorInfo           string
	Manufacturer         uint32
	FirmwareVersionMajor int
	FirmwareVersionMinor int
}

func (t *TPM) Info(ctx context.Context) (Info, error) {
	result := Info{}
	if err := t.Open(ctx); err != nil {
		return result, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer t.Close(ctx, false)

	a, err := attest.OpenTPM(t.attestConfig) // TODO: add layer of abstraction here, to ease testing?
	if err != nil {
		return result, fmt.Errorf("failed opening TPM: %w", err)
	}
	defer a.Close()

	info, err := a.Info()
	if err != nil {
		return result, fmt.Errorf("failed getting TPM info: %w", err)
	}

	result.FirmwareVersionMajor = info.FirmwareVersionMajor
	result.FirmwareVersionMinor = info.FirmwareVersionMinor
	result.Interface = uint8(info.Interface)
	result.Manufacturer = uint32(info.Manufacturer)
	result.VendorInfo = info.VendorInfo
	result.Version = uint8(info.Version)

	return result, nil
}
