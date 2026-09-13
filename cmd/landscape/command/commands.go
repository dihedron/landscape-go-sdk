package command

import (
	"github.com/dihedron/landscape-go-sdk/cmd/landscape/command/computer"
	"github.com/dihedron/landscape-go-sdk/internal/command/version"
)

// Commands is the set of root command groups.
type Commands struct {
	// Computer manages computers in Landscape.
	Computer computer.Computer `command:"computer" alias:"c" description:"Manage computers in Landscape."`
	// Version prints the program version information.
	Version version.Version `command:"version" alias:"v" description:"Print program version information."`
}
