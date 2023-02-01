package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	Version     string
	ReleaseDate string
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print the version information",
		Long:  `Print the version information.`,
		Run: func(cmd *cobra.Command, args []string) {
			if Version == "" {
				Version = "0000000-dev"
			}
			if ReleaseDate == "" {
				ReleaseDate = time.Now().UTC().Format("2006-01-02 15:04 MST")
			}

			fmt.Printf("%s/%s (%s/%s)\n", cmd.Parent().Name(), Version, runtime.GOOS, runtime.GOARCH)
			fmt.Printf("Release Date: %s\n", ReleaseDate)
		},
	}
}
