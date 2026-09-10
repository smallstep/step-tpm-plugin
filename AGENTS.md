# AGENTS.md

Guidance for AI coding agents working in this repository.

## Overview

`step-tpm-plugin` is a small Cobra CLI (module `github.com/smallstep/step-tpm-plugin`) for interacting with TPM 2.0 devices: reading TPM info and EKs, creating and managing AKs and keys, generating random bytes, and optionally running a TPM simulator. It is a `step` plugin: the `step` CLI runs any `step-<name>-plugin` binary found in `$(step path)/plugins` or on `PATH`, so this binary is invoked as `step tpm <command>`. All TPM logic lives in `go.step.sm/crypto/tpm` (the `smallstep/crypto` repo); this repo is the command-line surface over it. The project is marked beta in the README and its output formats may change.

## Commands

```bash
make build        # go build to bin/step-tpm-plugin (sets cmd.Version and cmd.ReleaseDate via ldflags)
make build-dev    # same, with -tags tpmsimulator (enables `simulator run`; needs OpenSSL headers, see README)
make test         # go test -coverprofile=coverage.out ./...
make fmt          # goimports -local github.com/smallstep/step-tpm-plugin
make lint         # golangci-lint (config curled from smallstep/workflows) + govulncheck
make generate     # go generate ./... then regenerate completions/* from the built binary
make bootstrap    # install golangci-lint, govulncheck, gotestsum
make install      # install bin/step-tpm-plugin under $(INSTALL_PREFIX)bin (default /usr/)
```

- `go build -o step-tpm-plugin -tags tpmsimulator .` is the README's simulator build. On macOS export `CGO_CFLAGS="-I$(brew --prefix openssl)/include"` and `CGO_LDFLAGS="-L$(brew --prefix openssl)/lib"` first.
- There are no `_test.go` files yet (the README lists tests as a TODO), so `make test` only compiles packages. Run a single test with `go test -run TestName ./path/to/package` once tests exist.
- CI (`.github/workflows/ci.yml`) calls the shared `smallstep/workflows` `goCI` workflow: lint, govulncheck, CodeQL, `gotestsum -- -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...`, and `V=1 make build`. It runs on pull requests; the `push` trigger names `master`, not `main`, so pushes to the default branch do not run it.
- `make lint` fetches `.golangci.yml` from the `master` branch of `smallstep/workflows`, which no longer exists; CI's lint job supplies its own config, so run `golangci-lint run` locally with an explicit config if needed.
- `make release` is a stub (`echo "TODO"`); there is no release workflow in this repo.
- Try the binary against hardware with `bin/step-tpm-plugin info` or `--device /path/to/socket info` against a running simulator (`sudo bin/step-tpm-plugin simulator run` after `make build-dev`).

## Generated Code - Do Not Edit

| Pattern | Generator |
|---------|-----------|
| `completions/{bash,fish,powershell,zsh}_completion` | `make generate` (Cobra `completion` subcommand of the built binary) |

There are no `go:generate` directives in the source; `go generate ./...` in `make generate` is a no-op today.

## Architecture

```
step-tpm-plugin/
├── main.go                  # calls cmd.Execute()
├── cmd/
│   ├── root.go              # rootCmd; step.Init() then registers subcommands
│   ├── info.go              # `info`: TPM version/manufacturer/firmware
│   ├── random.go            # `random`: TPM-generated random bytes (base64 or --hex)
│   ├── version.go           # `version`; Version/ReleaseDate set by ldflags
│   ├── ek.go, ek/get.go     # `ek get`: endorsement keys (table, --json, --pem, --all)
│   ├── ak.go, ak/*.go       # `ak create|list|get|delete`: attestation keys
│   ├── key.go, keys/*.go    # `key create|list|get|delete`; create takes --ak, --kty, --size
│   └── simulator.go, simulator/
│       ├── run.go               # `simulator run` flags: --socket, --seed, --verbose
│       ├── simulator_enabled.go # //go:build tpmsimulator: serves the simulator on a UNIX socket
│       └── simulator_disabled.go# default build: returns an error explaining the build tag
├── internal/
│   ├── command/             # command.New(usage, short, long, runner, preparers, finalizers)
│   │   ├── command.go       # runE: builds ctx, runs common + per-command preparers, runner, finalizers
│   │   ├── preparers.go     # RequireTPMWithStorage / RequireTPMWithoutStorage put a *tpm.TPM in ctx
│   │   └── context.go       # cobra.Command in/out of context
│   ├── flag/                # typed flag structs (Bool/String/Int), flag.Add, shared flag constructors
│   │   └── context.go       # flag.GetString/GetBool/GetInt/FirstArg read the pflag.FlagSet from ctx
│   └── render/render.go     # render.JSON: indented JSON to a writer
├── completions/             # generated shell completions
└── Makefile
```

### Command flow

1. `cmd.Execute()` runs `step.Init()` (from `cli-utils`) so `step.Path()` resolves `STEPPATH`, then `rootCmd.Execute()`.
2. Every command is built with `command.New(...)`. Its `RunE` stores the `*cobra.Command` and its `pflag.FlagSet` in `ctx`, then runs preparers in order: the common `fallbackTPMStore` (a no-op key store), then the command's own preparers.
3. `command.RequireTPMWithoutStorage` (info, ek, random) opens the TPM named by `--device`; `command.RequireTPMWithStorage` (ak, key) additionally wires a `storage.TPMStore`: a single-file store when `--storage-file` is set, otherwise a directory store at `--storage-directory` (default `$(step path)/tpm`). Both put the `*tpm.TPM` in `ctx` via `tpm.NewContext`.
4. The runner pulls values with `tpm.FromContext(ctx)` and `flag.GetX(ctx, name)` and prints a `go-pretty` table by default or JSON with `--json`.

## Conventions

- **CLI framework**: Cobra via the `internal/command` wrapper. Do not set `RunE` directly; pass a `command.Runner` and preparers to `command.New`. Subcommand groups (`ak`, `key`, `ek`, `simulator`) pass a nil runner.
- **Flags**: declare reusable flags as constructors in `internal/flag/flag.go` (`flag.Device()`, `flag.JSON()`, ...) and name constants as `flag.FlagXxx`; one-off flags use `flag.Int{...}` / `flag.String{...}` literals inline. Read them with `flag.GetString(ctx, ...)` etc.; these panic if the flag was not added to the command, so add before reading.
- **Errors**: `fmt.Errorf("...: %w", err)` with a short verb phrase (`"creating key failed: %w"`). No `pkg/errors`.
- **Logging**: none outside the simulator, which uses stdlib `log` gated by `--verbose` through `logF`.
- **Output**: table via `github.com/jedib0t/go-pretty/table` mirrored to stdout; `--json` via `render.JSON(os.Stdout, v)`; `--pem` where certificates or keys are involved.
- **Build tags**: simulator support is compiled only with `tpmsimulator`; keep the enabled/disabled file pair in sync when changing the `runSimulator` signature.
- **Testing**: no tests yet. If adding some, prefer `go.step.sm/crypto/tpm/simulator` for TPM-backed tests (it needs the `tpmsimulator` tag and CGO/OpenSSL).

## Environment variables

- `STEP_TPM_DEVICE` - default value for `--device`.
- `STEPPATH` - step configuration root, read by `cli-utils`; the default key storage directory is `$STEPPATH/tpm`.

## Related repos

- `smallstep/crypto` (`go.step.sm/crypto`): `tpm`, `tpm/storage`, `tpm/simulator`, `pemutil`, `randutil`. Most feature work touches this module first; test against a local checkout with `replace go.step.sm/crypto => ../crypto` in `go.mod` (do not commit the replace).
- `smallstep/cli-utils`: `step.Init()` / `step.Path()`.
- `smallstep/cli`: the `step` binary that discovers and executes this plugin.
- `google/go-tpm-tools`: provides the simulator (via `crypto`); its README covers OpenSSL build errors.
