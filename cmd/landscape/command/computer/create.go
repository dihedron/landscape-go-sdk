package computer

/*
type Create struct {
	base.Command
	User *midpoint.User `short:"d" long:"data" description:"New user's data, either as an inline value or as a @file (in JSON or YAML format)."`
}

func (cmd *Create) Execute(args []string) error {
	slog.Debug("running user create command", "endpoint", cmd.Endpoint, "username", cmd.Username, "password", cmd.Password, "user", *cmd.User)

	mp := cmd.GetAPI()
	defer mp.Close()

	id, err := mp.User.Create(context.Background(), cmd.User)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return err
	}
	slog.Debug("user created", "id", id)
	cmd.Write(os.Stdout, midpoint.ID{ID: id})
	return nil
}
*/
