//go:build tpmsimulator

package simulator

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/smallstep/panoramix/v5/logware"
	"go.step.sm/crypto/randutil"
	"go.step.sm/crypto/tpm"
	"go.step.sm/crypto/tpm/simulator"

	"github.com/smallstep/step-tpm-plugin/internal/flag"
)

const (
	loggerName = "tpm"
)

func runSimulator(ctx context.Context) (err error) {
	var (
		socket  = flag.GetString(ctx, flag.FlagSocket)
		seed    = flag.GetString(ctx, flag.FlagSeed)
		verbose = flag.GetBool(ctx, flag.FlagVerbose)
		logger  *slog.Logger
	)

	loggerOptions := []logware.Option{logware.WithName(loggerName)}
	if verbose {
		loggerOptions = append(loggerOptions, logware.WithLevel(slog.LevelDebug))
	}

	logger = logware.Logger(loggerOptions...)

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

	if err := printTPMInfo(info, eks, socket, seed); err != nil {
		logger.ErrorContext(ctx, "failed printing TPM info", logware.Error(err))
	}

	logger.InfoContext(ctx, "TPM simulator available", slog.String("socket", socket))

	// register shutdown signals and perform cleanup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.InfoContext(ctx, "stopping TPM simulator ...", slog.String("socket", socket))
		os.Remove(socket)
		os.Exit(0)
	}()

	for {
		// accept incoming connections
		conn, err := ln.Accept()

		// TODO: perform an early, internal connection test?

		// check if connection was established successfully
		if err != nil {
			logger.ErrorContext(ctx, "connection error", logware.Error(err))
		} else {
			logger.DebugContext(ctx, "connection established", slog.String("addr", conn.RemoteAddr().String()))
		}

		// handle the connection in a separate goroutine.
		go func(conn net.Conn) {
			defer conn.Close()

			// create buffer and read data from the new connection
			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			if err != nil {
				logger.ErrorContext(ctx, "failed reading fom socket", logware.Error(err))
			}

			logger.DebugContext(ctx, "read from socket", slog.Any("bytes", buf[:n]))

			nw, err := s.Write(buf[:n])
			if err != nil {
				logger.ErrorContext(ctx, "failed writing to TPM", logware.Error(err))
			}

			logger.DebugContext(ctx, "written to TPM", slog.Any("bytes", nw))

			newBuf := make([]byte, 4096)
			nr, err := s.Read(newBuf)
			if err != nil {
				logger.ErrorContext(ctx, "failed reading from TPM", logware.Error(err))
			}

			logger.DebugContext(ctx, "read from TPM", slog.Any("bytes", nr))

			_, err = conn.Write(newBuf[:nr])
			if err != nil {
				logger.ErrorContext(ctx, "failed writing to socket", logware.Error(err))
			}

			logger.DebugContext(ctx, "written to socket", slog.Any("bytes", newBuf[:nr]))

			// TODO(hs): log at least something about the interaction? The simulator is now very quiet after starting
		}(conn)
	}
}

var (
	purple       = lipgloss.Color("99")
	gray         = lipgloss.Color("245")
	lightGray    = lipgloss.Color("241")
	headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	cellStyle    = lipgloss.NewStyle().Padding(0, 1).MaxWidth(80)
	oddRowStyle  = cellStyle.Foreground(gray)
	evenRowStyle = cellStyle.Foreground(lightGray)
)

func printTPMInfo(info *tpm.Info, eks []*tpm.EK, socket, seed string) error {
	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			case row%2 == 0:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		})

	tbl.Row("Version", info.Version.String())
	tbl.Row("Interface", fmt.Sprintf("%s (simulator)", info.Interface))
	tbl.Row("Manufacturer", info.Manufacturer.String())
	tbl.Row("Vendor Info", info.VendorInfo)
	tbl.Row("Firmware Version", info.FirmwareVersion.String())
	for _, ek := range eks {
		u, err := ek.FingerprintURI()
		if err != nil {
			return err
		}
		tbl.Row(fmt.Sprintf("EK URI (%s)", ek.Type()), u.String())
	}
	tbl.Row("UNIX socket", socket)
	tbl.Row("Seed", seed)

	w := colorprofile.NewWriter(os.Stdout, os.Environ())
	fmt.Fprintln(w, tbl)

	return nil
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
