package ek

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"go.step.sm/crypto/tpm"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/internal/render"
)

func NewGetEKCommand() *cobra.Command {
	const (
		long = `print TPM EK details.
`
		short = "print TPM EK details"
	)

	cmd := command.New("get", short, long, runGetEK, []command.Preparer{command.RequireTPMWithoutStorage}, nil)

	flag.Add(cmd,
		flag.JSON(),
		flag.Device(),
		flag.PEM(),
		flag.Bool{
			Name:        "all",
			Description: "Print all availables EKs",
		},
	)

	return cmd
}

func runGetEK(ctx context.Context) error {
	var (
		t    = tpm.FromContext(ctx)
		json = flag.GetBool(ctx, flag.FlagJSON)
		pem  = flag.GetBool(ctx, flag.FlagPEM)
		all  = flag.GetBool(ctx, "all")
	)

	eks, err := t.GetEKs(ctx)
	if err != nil {
		return fmt.Errorf("error getting EKs: %w", err)
	}

	switch {
	case pem:
		if all {
			for _, ek := range eks {
				b, err := ek.PEM()
				if err != nil {
					return err
				}
				fmt.Println(b)
			}
			return nil
		}

		b, err := eks[0].PEM()
		if err != nil {
			return err
		}
		fmt.Println(b)

	case json:
		if all {
			return render.JSON(os.Stdout, eks)
		}

		return render.JSON(os.Stdout, eks[0])
	default:
		t1 := table.NewWriter()
		t1.SetOutputMirror(os.Stdout)
		t1.AppendHeader(table.Row{"Type", "Certificate", "CertificateURL"})
		for _, ek := range eks {
			cert := "-"
			if ek.Certificate() != nil {
				cert = "OK"
			}

			t1.AppendRow(table.Row{ek.Type(), cert, ek.CertificateURL()})
		}
		t1.Render()
	}

	return nil
}
