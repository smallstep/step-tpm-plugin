package ak

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/pkg/tpm"
)

// TODO(hs): move to internal?

func NewListAKCommand() *cobra.Command {
	const (
		long = `loooooooong list description.
`
		short = "short list description"
	)

	cmd := command.New("list", short, long, runListAK, []command.Preparer{command.RequireTPM}, nil)

	flag.Add(cmd,
		flag.StorageFile(),
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

	for _, ak := range aks {
		fmt.Println(ak.Name, ak.Data)
	}

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
