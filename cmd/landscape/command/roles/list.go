package roles

import (
	"log/slog"
	"os"

	"github.com/dihedron/midpoint-go-sdk/pkg/metadata"
)

// List is the command that lists all roles.
type List struct {
	// Verbose is the flag that indicates whether to print verbose information about the application.
	Verbose bool `short:"v" long:"verbose" description:"Print verbose information about the application."`
}

// Execute is the real implementation of the Version command.
func (cmd *List) Execute(args []string) error {
	slog.Debug("running version command")
	if cmd.Verbose {
		metadata.PrintFull(os.Stdout)
	} else {
		metadata.Print(os.Stdout)
	}
	slog.Debug("command done")
	return nil
}
