package ek

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm"
)

func NewGetEKCommand() *cobra.Command {
	const (
		long = `print TPM EK details.
`
		short = "print TPM EK details"
	)

	cmd := command.New("get", short, long, runGetEK, []command.Preparer{command.RequireTPM}, nil)

	flag.Add(cmd,
		flag.JSON(),
		flag.Device(),
		flag.StorageFile(),
	)

	return cmd
}

func runGetEK(ctx context.Context) error {

	var (
		t    = tpm.FromContext(ctx)
		json = flag.GetBool(ctx, flag.FlagJSON)
	)

	_ = json

	eks, err := t.GetEKs(ctx)
	if err != nil {
		return fmt.Errorf("error getting EKs: %w", err)
	}

	t1 := table.NewWriter()
	t1.SetOutputMirror(os.Stdout)
	t1.AppendHeader(table.Row{"Public Key", "Certificate", "CertificateURL"})
	for _, ek := range eks {
		t1.AppendRow(table.Row{fmt.Sprintf("%T", ek.Public), ek.Certificate, ek.CertificateURL})
	}
	t1.Render()

	return nil
}
