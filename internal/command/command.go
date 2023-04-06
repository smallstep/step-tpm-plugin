package command

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/smallstep/step-tpm-plugin/internal/flag"
)

type (
	Preparer  func(context.Context) (context.Context, error)
	Runner    func(context.Context) error
	Finalizer func(context.Context) (context.Context, error)
)

func New(usage, short, long string, fn Runner, p []Preparer, f []Finalizer) *cobra.Command {
	return &cobra.Command{
		Use:   usage,
		Short: short,
		Long:  long,
		RunE:  runE(fn, p, f...),
	}
}

var commonPreparers = []Preparer{
	fallbackTPMStore,
}

var commonFinalizers = []Finalizer{}

func runE(fn Runner, preparers []Preparer, finalizers ...Finalizer) func(*cobra.Command, []string) error {
	if fn == nil {
		return nil
	}

	return func(cmd *cobra.Command, _ []string) (err error) {
		ctx := cmd.Context()
		ctx = NewContext(ctx, cmd)
		ctx = flag.NewContext(ctx, cmd.Flags())

		// run the common preparers
		if ctx, err = prepare(ctx, commonPreparers...); err != nil {
			return
		}

		// run the preparers specific to the command
		if ctx, err = prepare(ctx, preparers...); err != nil {
			return
		}

		// run the command
		if err = fn(ctx); err == nil {
			// run the finalizers registered for the command
			if ctx, err = finalize(ctx, finalizers...); err != nil {
				return err
			}

			// and finally, run common finalizers
			if _, err = finalize(ctx, commonFinalizers...); err != nil {
				return err
			}
		}

		return
	}
}

func prepare(parent context.Context, preparers ...Preparer) (ctx context.Context, err error) {
	ctx = parent
	for _, p := range preparers {
		if ctx, err = p(ctx); err != nil {
			break
		}
	}

	return
}

func finalize(parent context.Context, finalizers ...Finalizer) (ctx context.Context, err error) {
	ctx = parent
	for _, f := range finalizers {
		if ctx, err = f(ctx); err != nil {
			break
		}
	}

	return
}
