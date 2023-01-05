package flag

import "github.com/spf13/cobra"

const (
	FlagDeviceName  = "device"
	FlagStorageFile = "storage-file"
	FlagJSON        = "json"
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
	}
}

// StorageFile returns a "storage-file" string flag.
func StorageFile() String {
	return String{
		Name:        FlagStorageFile,
		Shorthand:   "s",
		Description: "Filename for TPM key storage",
		Default:     "data.json",
	}
}
