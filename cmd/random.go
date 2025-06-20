package cmd

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"

	"github.com/spf13/cobra"

	"go.step.sm/crypto/tpm"

	"github.com/smallstep/step-tpm-plugin/internal/command"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
)

func NewRandomCommand() *cobra.Command {
	const (
		long = `subcommand for generating random data.
`
		short = "subcommand for generating random data"
	)

	cmd := command.New("random <command>", short, long, runRandom, []command.Preparer{command.RequireTPMWithoutStorage}, nil)

	flag.Add(cmd,
		flag.Device(),
		flag.Int{
			Name:        "size",
			Description: "number of random bytes to generate",
			Default:     32,
		},
		flag.Bool{
			Name:        "hex",
			Description: "output using hexadecimal characters",
			Default:     false,
		},
	)

	return cmd
}

func runRandom(ctx context.Context) error {
	var (
		t         = tpm.FromContext(ctx)
		size      = flag.GetInt(ctx, "size")
		outputHex = flag.GetBool(ctx, "hex")
	)

	// while the underlying library takes uint16 as the max, the actual maximum that a TPM will
	// return is likely a lot lower. I've observed 64 bytes being returned, but it might as well be
	// 32 bytes. Does that warrant changing the max below?
	if size < 0 || size > math.MaxUint16 {
		return fmt.Errorf("'--size' must be between 0 and %d; got %d", math.MaxUint16, size)
	}

	r, err := t.GenerateRandom(ctx, uint16(size))
	if err != nil {
		return fmt.Errorf("failed generating random: %w", err)
	}

	var out string
	if outputHex {
		out = hex.EncodeToString(r)
	} else {
		out = base64.StdEncoding.EncodeToString(r)
	}

	fmt.Println(out)

	return nil
}
