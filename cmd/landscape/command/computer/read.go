package computer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/dihedron/landscape-go-sdk/internal/command/base"
	"github.com/dihedron/landscape-go-sdk/pkg/landscape"
)

type Read struct {
	base.Command
}

func (cmd *Read) Execute(args []string) error {
	slog.Debug("running computer read command", "endpoint", cmd.Endpoint, "email", cmd.Email, "password", cmd.Password, "account", cmd.Account, "ids", args)
	if len(args) == 0 {
		slog.Error("no ids provided")
		return fmt.Errorf("at least one ID must be provided")
	}

	api := cmd.GetAPI()
	defer api.Close()

	var result error
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			slog.Error("error parsing computer id", "id", arg, "error", err)
			errors.Join(result, err)
			continue
		}
		slog.Debug("reading computer", "id", arg)
		self, err := api.Computer.Read(context.Background(), id, landscape.ComputerReadOptions{})
		if err != nil {
			slog.Error("error reading computer", "id", arg, "error", err)
			errors.Join(result, err)
			continue
		}
		if err = cmd.Write(os.Stdout, self); err != nil {
			errors.Join(result, err)
		}
	}
	return result
}
