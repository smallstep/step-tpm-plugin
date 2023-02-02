package csr

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/pkg/tpm"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm/skae"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
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

	fmt.Println(signer.Public())

	params, err := key.CertificationParameters(ctx)
	if err != nil {
		return fmt.Errorf("error getting key certification parameters: %w", err)
	}

	// TODO: fix the SKAE. It fails with "creating SKAE extension failed: creating SKAE extension failed: asn1: structure error: invalid object identifier"
	skaeExtension, err := skae.CreateSubjectKeyAttestationEvidenceExtension(nil, params)
	if err != nil {
		return fmt.Errorf("creating SKAE extension failed: %w", err)
	}

	_ = skaeExtension

	// TODO: alter existing CSR instead; take that as argument and/or move this to library?
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: "test cn",
		},
		// ExtraExtensions: []pkix.Extension{
		// 	skaeExtension,
		// },
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, template, signer)
	if err != nil {
		return fmt.Errorf("error creating certificate request: %w", err)
	}

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return fmt.Errorf("error parsing certificate request: %w", err)
	}

	fmt.Printf("%#+v\n", csr)

	fmt.Println(base64.StdEncoding.EncodeToString(csr.Raw))

	// TODO: output to file? Output the CSR PEM formatted to stdout?

	return nil
}
