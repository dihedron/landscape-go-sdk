package base

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"

	"github.com/dihedron/landscape-go-sdk/pkg/landscape"
	"go.yaml.in/yaml/v3"
)

// Command is the base command.
type Command struct {
	// Endpoint specifies the network endpoint for Landscape API access.
	Endpoint string `short:"E" long:"endpoint" description:"The network endpoint for Landscape API access." required:"true" env:"LANDSCAPE_ENDPOINT"`
	// Email specifies the email for authentication.
	Email string `short:"U" long:"email" description:"The email for authentication." required:"true" env:"LANDSCAPE_EMAIL"`
	// Password specifies the password for authentication.
	Password string `short:"P" long:"password" description:"The password for authentication." required:"true" env:"LANDSCAPE_PASSWORD"`
	// Account specifies the account for authentication.
	Account string `short:"A" long:"account" description:"The account for authentication." optional:"true" env:"LANDSCAPE_ACCOUNT"`
	// Format specifies the output format.
	//lint:ignore SA5008 duplicate alias tags are legitimate
	Format string `short:"F" long:"format" description:"The format of the output." optional:"true" default:"yaml" choice:"text" choice:"json" choice:"yaml" choice:"none" env:"MIDPOINT_FORMAT"`
	// Debug enables debug mode.
	Debug bool `short:"D" long:"debug" description:"Enable debug mode." optional:"true" env:"MIDPOINT_DEBUG"`
	// Enable saving response to file.
	ResponseSavePath *string `long:"save-response-to-file" description:"Enable saving HTTP responses to the given file." optional:"true" env:"MIDPOINT_RESPONSE_SAVE_TO_FILE"`
}

func (cmd *Command) GetAPI() *landscape.API {
	options := []landscape.Option{
		landscape.WithBasicAuth(cmd.Email, cmd.Password, cmd.Account),
		landscape.WithDebug(cmd.Debug),
		landscape.WithTraceRequest(true),
	}
	if cmd.ResponseSavePath != nil {
		options = append(options, landscape.WithSaveResponse(true, *cmd.ResponseSavePath))
	}

	return landscape.New(cmd.Endpoint, options...)
}

func (cmd *Command) Write(stream io.Writer, object any) error {
	switch cmd.Format {
	case "yaml":
		data, err := yaml.Marshal(object)
		if err != nil {
			return err
		}
		fmt.Fprintf(stream, "%s", string(data))
	case "json":
		data, err := json.Marshal(
			object,
			json.OmitZeroStructFields(true),
			jsontext.WithIndent("  "),
		)
		if err != nil {
			return err
		}
		fmt.Fprintf(stream, "%s\n", string(data))
	case "text":
		fmt.Fprintf(stream, "%v\n", object)
	default:
		return nil
	}
	return nil
}
