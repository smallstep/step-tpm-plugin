package cmd

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

func NewInfoCommand() *cobra.Command {
	const (
		long = `retrieve TPM information.
`
		short = "retrieve TPM information"
	)

	cmd := command.New("info", short, long, runInfo, []command.Preparer{command.RequireTPMWithoutStorage}, nil)

	cmd.Args = cobra.NoArgs

	flag.Add(cmd,
		flag.Device(),
		flag.JSON(),
	)

	return cmd
}

func runInfo(ctx context.Context) error {

	var (
		t    = tpm.FromContext(ctx)
		json = flag.GetBool(ctx, flag.FlagJSON)
	)

	info, err := t.Info(ctx)
	if err != nil {
		return fmt.Errorf("failed getting TPM info: %w", err)
	}

	// TODO(hs): what if there are multiple TPMs?
	// TODO(hs): output additional data with --verbose?
	// TODO(hs): add basic EK info to output?

	_ = json // TODO(hs): actual JSON output

	t1 := table.NewWriter()
	t1.SetOutputMirror(os.Stdout)
	t1.AppendRows([]table.Row{
		{"Version", info.Version},
		{"Interface", info.Interface},
		{"Manufacturer", info.Manufacturer},
		{"Vendor Info", info.VendorInfo},
		{"Firmware Version", fmt.Sprintf("%d.%d", info.FirmwareVersionMajor, info.FirmwareVersionMinor)},
	})
	t1.Render()

	return nil
}
