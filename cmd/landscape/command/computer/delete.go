package computer

// type Delete struct {
// 	base.Command
// }

// func (cmd *Delete) Execute(args []string) error {
// 	slog.Debug("running user delete command", "endpoint", cmd.Endpoint, "username", cmd.Username, "password", cmd.Password, "ids", args)
// 	if len(args) == 0 {
// 		slog.Error("no ids provided")
// 		return fmt.Errorf("at least one ID must be provided")
// 	}

// 	mp := cmd.GetAPI()
// 	defer mp.Close()

// 	var result error
// 	for _, arg := range args {
// 		slog.Debug("deleting user", "id", arg)
// 		err := mp.User.Delete(context.Background(), arg)
// 		if err != nil {
// 			slog.Error("error deleting user", "id", arg, "error", err)
// 			errors.Join(result, err)
// 			continue
// 		}
// 		if err = cmd.Write(os.Stdout, arg); err != nil {
// 			errors.Join(result, err)
// 		}
// 	}
// 	return result
// }
