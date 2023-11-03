//go:build !tpmsimulator
// +build !tpmsimulator

package simulator

import (
	"context"
	"errors"
)

func runSimulator(context.Context) error {
	return errors.New(`simulator is only available when built with the "tpmsimulator" tag`)
}
