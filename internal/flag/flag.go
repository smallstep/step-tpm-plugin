package flag

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"go.step.sm/cli-utils/step"
)

const (
	FlagDeviceName       = "device"
	FlagStorageFile      = "storage-file"
	FlagStorageDirectory = "storage-directory"
	FlagJSON             = "json"
	FlagPEM              = "pem"
	FlagBundle           = "bundle"
	FlagAK               = "ak"
	FlagBlob             = "blob"
	FlagPrivate          = "private"
	FlagPublic           = "public"
	FlagTSS2             = "tss2"
	FlagSocket           = "socket"
	FlagSeed             = "seed"
	FlagVerbose          = "verbose"
)

// Flag wraps the set of flags.
type Flag interface {
	addTo(*cobra.Command)
}

type Set []Flag

func (s Set) addTo(cmd *cobra.Command) {
	for _, flag := range s {
		flag.addTo(cmd)
	}
}

// Add adds flag to cmd, binding them on v should v not be nil.
func Add(cmd *cobra.Command, flags ...Flag) {
	for _, flag := range flags {
		flag.addTo(cmd)
	}
}

// Bool wraps the set of boolean flags.
type Bool struct {
	Name        string
	Shorthand   string
	Description string
	Default     bool
	Hidden      bool
}

func (b Bool) addTo(cmd *cobra.Command) {
	flags := cmd.Flags()

	if b.Shorthand != "" {
		_ = flags.BoolP(b.Name, b.Shorthand, b.Default, b.Description)
	} else {
		_ = flags.Bool(b.Name, b.Default, b.Description)
	}

	f := flags.Lookup(b.Name)
	f.Hidden = b.Hidden
}

// String wraps the set of string flags.
type String struct {
	Name        string
	Shorthand   string
	Description string
	Default     string
	ConfName    string
	EnvName     string
	Hidden      bool
}

func (s String) addTo(cmd *cobra.Command) {
	flags := cmd.Flags()

	if s.Shorthand != "" {
		_ = flags.StringP(s.Name, s.Shorthand, s.Default, s.Description)
	} else {
		_ = flags.String(s.Name, s.Default, s.Description)
	}

	f := flags.Lookup(s.Name)
	f.Hidden = s.Hidden
}

// Int wraps the set of int flags.
type Int struct {
	Name        string
	Shorthand   string
	Description string
	Default     int
	Hidden      bool
}

func (i Int) addTo(cmd *cobra.Command) {
	flags := cmd.Flags()

	if i.Shorthand != "" {
		_ = flags.IntP(i.Name, i.Shorthand, i.Default, i.Description)
	} else {
		_ = flags.Int(i.Name, i.Default, i.Description)
	}

	f := flags.Lookup(i.Name)
	f.Hidden = i.Hidden
}

// JSON returns a "json" bool flag.
func JSON() Bool {
	return Bool{
		Name:        FlagJSON,
		Description: "Output in JSON format",
	}
}

// Device returns a "device" string flag.
func Device() String {
	return String{
		Name:        FlagDeviceName,
		Shorthand:   "d",
		Description: "TPM device name",
		Default:     os.Getenv("STEP_TPM_DEVICE"),
	}
}

// StorageFile returns a "storage-file" string flag.
func StorageFile() String {
	return String{
		Name:        FlagStorageFile,
		Description: "Filename for TPM key storage",
		Default:     "",
	}
}

// StorageDirectory returns a "storage-directory" string flag.
func StorageDirectory() String {
	return String{
		Name:        FlagStorageDirectory,
		Shorthand:   "s",
		Description: "Directory to store TPM keys",
		Default:     filepath.Join(step.Path(), "tpm"),
	}
}

// PEM returns a "pem" bool flag.
func PEM() Bool {
	return Bool{
		Name:        FlagPEM,
		Description: "Output in PEM format",
	}
}

// Bundle returns a "bundle" bool flag
func Bundle() Bool {
	return Bool{
		Name:        FlagBundle,
		Description: "Output certificate chain bundle",
	}
}

// AK returns an "ak" string flag.
func AK() String {
	return String{
		Name:        FlagAK,
		Description: "Name of the AK to attest new key with",
	}
}

// Blob returns a "blob" bool flag
func Blob() Bool {
	return Bool{
		Name:        FlagBlob,
		Description: "Print a blob",
	}
}

// Private returns a "private" bool flag
func Private() Bool {
	return Bool{
		Name:        FlagPrivate,
		Description: "Print private blob",
	}
}

// Public returns a "public" bool flag
func Public() Bool {
	return Bool{
		Name:        FlagPublic,
		Description: "Print public blob",
	}
}

// TSS2 returns a "tts2" bool flag
func TSS2() Bool {
	return Bool{
		Name:        FlagTSS2,
		Description: "Print the public and private blobs using the TSS2 format",
	}
}

// Socket returns a "socket" string flag.
func Socket() String {
	return String{
		Name:        FlagSocket,
		Description: "Path to UNIX socket to serve TPM simulator on",
	}
}

// Seed returns a "seed" string flag
func Seed() String {
	return String{
		Name:        FlagSeed,
		Description: "Seed value for TPM simulator",
	}
}

// Verbose returns a "verbose" bool flag
func Verbose() Bool {
	return Bool{
		Name:        FlagVerbose,
		Description: "Enable verbose logging",
	}
}
