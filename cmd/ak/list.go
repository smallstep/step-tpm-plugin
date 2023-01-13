package ak

import (
	"context"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm"
)

// TODO(hs): move to internal?

func NewListAKCommand() *cobra.Command {
	const (
		long = `list AK details.
`
		short = "list AK details"
	)

	cmd := command.New("list", short, long, runListAK, []command.Preparer{command.RequireTPMWithStorage}, nil)

	flag.Add(cmd,
		flag.StorageFile(),
		flag.StorageDirectory(),
		flag.Device(),
		flag.JSON(),
	)

	return cmd
}

func runListAK(ctx context.Context) error {

	var (
		t = tpm.FromContext(ctx)
	)

	aks, err := t.ListAKs(ctx)
	if err != nil {
		return err
	}

	t1 := table.NewWriter()
	t1.SetOutputMirror(os.Stdout)
	t1.AppendHeader(table.Row{"Name", "Data length"})
	for _, ak := range aks {
		t1.AppendRow(table.Row{ak.Name, len(ak.Data)})
	}
	t1.Render()

	// cfg := config.FromContext(ctx)
	// client := client.FromContext(ctx)

	// var apps []api.App
	// if apps, err = client.API().GetApps(ctx, nil); err != nil {
	// 	return
	// }

	// out := iostreams.FromContext(ctx).Out
	// if cfg.JSONOutput {
	// 	_ = render.JSON(out, apps)

	// 	return
	// }

	// rows := make([][]string, 0, len(apps))
	// for _, app := range apps {
	// 	latestDeploy := ""
	// 	if app.Deployed && app.CurrentRelease != nil {
	// 		latestDeploy = format.RelativeTime(app.CurrentRelease.CreatedAt)
	// 	}

	// 	rows = append(rows, []string{
	// 		app.Name,
	// 		app.Organization.Slug,
	// 		app.Status,
	// 		app.PlatformVersion,
	// 		latestDeploy,
	// 	})
	// }

	// _ = render.Table(out, "", rows, "Name", "Owner", "Status", "Platform", "Latest Deploy")

	// return

	return nil
}
