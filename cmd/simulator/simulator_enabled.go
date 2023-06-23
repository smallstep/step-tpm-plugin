//go:build tpmsimulator
// +build tpmsimulator

package simulator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"go.step.sm/crypto/randutil"
	"go.step.sm/crypto/tpm"

	"github.com/jedib0t/go-pretty/table"
	"github.com/smallstep/step-tpm-plugin/internal/flag"
	"github.com/smallstep/step-tpm-plugin/internal/simulator"
)

func runSimulator(ctx context.Context) (err error) {
	var (
		socket  = flag.GetString(ctx, flag.FlagSocket)
		seed    = flag.GetString(ctx, flag.FlagSeed)
		verbose = flag.GetBool(ctx, flag.FlagVerbose)
	)

	if seed == "" {
		if seed, err = randutil.Hex(16); err != nil {
			return fmt.Errorf("failed generating TPM seed: %w", err)
		}
	}

	if socket == "" {
		if socket, err = getTPMSimulatorSocketPath(); err != nil {
			return
		}
	}

	s, err := simulator.New(simulator.WithSeed(seed))
	if err != nil {
		return
	}

	err = s.Open()
	if err != nil {
		return
	}
	defer func() {
		closeErr := s.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	// create a new TPM instance backed by simulator to
	// test that the simulator is functioning correctly.
	t, err := tpm.New(tpm.WithSimulator(s))
	if err != nil {
		return
	}

	info, err := t.Info(ctx) // TODO(hs): validate expected properties?
	if err != nil {
		return
	}

	eks, err := t.GetEKs(ctx)
	if err != nil {
		return
	}

	ln, err := net.Listen("unix", socket)
	if err != nil {
		return
	}
	defer func() {
		removeErr := os.Remove(socket)
		if removeErr != nil && err == nil {
			err = removeErr
		}
	}()

	t1 := table.NewWriter()
	t1.SetOutputMirror(os.Stdout)
	t1.AppendRows([]table.Row{
		{"Version", info.Version},
		{"Interface", fmt.Sprintf("%s (simulator)", info.Interface)},
		{"Manufacturer", info.Manufacturer},
		{"Vendor Info", info.VendorInfo},
		{"Firmware Version", info.FirmwareVersion},
	})
	for _, ek := range eks {
		u, err := ek.FingerprintURI()
		if err != nil {
			return err
		}
		t1.AppendRow(table.Row{
			fmt.Sprintf("EK URI (%s)", ek.Type()), u.String(),
		})
	}
	t1.AppendRows([]table.Row{
		{"UNIX socket", socket},
		{"Seed", seed},
	})
	t1.Render()

	logF(true, "TPM simulator available at %q\n", socket)

	// register shutdown signals and perform cleanup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logF(true, "stopping TPM simulator at %q\n", socket)
		os.Remove(socket)
		os.Exit(0)
	}()

	for {
		// accept incoming connections
		conn, err := ln.Accept()

		// TODO: perform an early, internal connection test?

		// check if connection was established successfully
		if err != nil {
			logF(true, "connection error: %v\n", err)
		} else {
			logF(verbose, "connection established: %v\n", conn)
		}

		// handle the connection in a separate goroutine.
		go func(conn net.Conn) {
			defer conn.Close()

			// create buffer and read data from the new connection
			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			if err != nil {
				logF(true, "failed reading from socket: %v", err)
			}

			logF(verbose, "read from socket: %v", buf[:n])

			nw, err := s.Write(buf[:n])
			if err != nil {
				logF(true, "failed writing to TPM: %v", err)
			}

			logF(verbose, "written to TPM: %v", nw)

			newBuf := make([]byte, 4096)
			nr, err := s.Read(newBuf)
			if err != nil {
				logF(true, "failed reading from TPM: %v", err)
			}

			logF(verbose, "read from TPM: %v", nr)

			_, err = conn.Write(newBuf[:nr])
			if err != nil {
				logF(true, "failed writing to socket: %v", err)
			}

			logF(verbose, "written to socket: %v", newBuf[:nr])

			// TODO(hs): log at least something about the interaction? The simulator is now very quiet after starting
		}(conn)
	}
}

func getTPMSimulatorSocketPath() (sockAddr string, err error) {
	paths := []string{"/run", "/var/run"}
	for _, dir := range paths {
		if _, err = os.Stat(dir); err == nil {
			sockAddr = filepath.Join(dir, "step-tpmsimulator.sock")
			return
		}
	}
	return "", errors.New("could not automatically determine TPM simulator socket path")
}

func logF(shouldPrint bool, format string, v ...any) {
	if shouldPrint {
		log.Printf(format, v...)
	}
}
