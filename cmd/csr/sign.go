package csr

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/google/go-attestation/attest"
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"go.step.sm/crypto/tpm"
	"go.step.sm/crypto/tpm/skae"
)

func NewSignCSRCommand() *cobra.Command {
	const (
		long = `temporary CSR sign command.
`
		short = "temporary CSR sign command"
	)

	cmd := command.New("sign", short, long, runSignCSR, []command.Preparer{command.RequireTPMWithStorage}, nil)

	cmd.Args = cobra.ExactArgs(1)

	flag.Add(cmd,
		flag.StorageFile(),
		flag.StorageDirectory(),
		flag.JSON(),
		flag.Device(),
	)

	return cmd
}

// TODO: remove me; temporary function until we get the abstraction right
func runSignCSR(ctx context.Context) error {
	var (
		t    = tpm.FromContext(ctx)
		name = flag.FirstArg(ctx)
	)

	key, err := t.GetKey(ctx, name)
	if err != nil {
		return fmt.Errorf("getting key failed: %w", err)
	}

	signer, err := key.Signer(ctx)
	if err != nil {
		return fmt.Errorf("getting signer for key failed: %w", err)
	}

	params, err := key.CertificationParameters(ctx)
	if err != nil {
		return fmt.Errorf("error getting key certification parameters: %w", err)
	}

	if !key.WasAttested() {
		return fmt.Errorf("key %q was not attested by an AK", key.Name())
	}

	ak, err := t.GetAK(ctx, key.AttestedBy())
	if err != nil {
		return fmt.Errorf("getting key failed: %w", err)
	}

	_ = ak

	// TODO: the AK cert should belong to an actual AK and thus be supplied
	// as a parameter, or automatically retrieved from the AK.
	fakeAK := &x509.Certificate{
		Issuer: pkix.Name{
			CommonName: "Test AK",
		},
		SerialNumber: big.NewInt(12345678),
		OCSPServer: []string{
			"https://www.example.com/ocsp/1",
			"https://www.example.com/ocsp/2",
		},
		IssuingCertificateURL: []string{
			"https://www.example.com/issuing/cert1",
		},
	}

	// retrieve AK attestation params and verify the key was attested by the AK
	attestParams, err := ak.AttestationParameters(ctx)
	if err != nil {
		return fmt.Errorf("getting AK attestation parameters failed: %w", err)
	}

	akPub, err := attest.ParseAKPublic(attest.TPMVersion20, attestParams.Public)
	if err != nil {
		return fmt.Errorf("failed to parse AK public: %w", err)
	}

	if err := params.Verify(attest.VerifyOpts{Public: akPub.Public, Hash: akPub.Hash}); err != nil {
		return fmt.Errorf("error verifying AK attested key: %w", err)
	}

	// TODO: determine whether or not to encrypt the evidence
	shouldEncrypt := false
	skaeExtension, err := skae.CreateSubjectKeyAttestationEvidenceExtension(fakeAK, params, shouldEncrypt)
	if err != nil {
		return fmt.Errorf("creating SKAE extension failed: %w", err)
	}

	// TODO: alter existing CSR instead; take that as argument and/or move this to library?
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: "test cn",
		},
		ExtraExtensions: []pkix.Extension{
			skaeExtension,
		},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, template, signer)
	if err != nil {
		return fmt.Errorf("error creating certificate request: %w", err)
	}

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return fmt.Errorf("error parsing certificate request: %w", err)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(csr.Raw))

	// TODO: output to file? Output the CSR PEM formatted to stdout?

	return nil
}
