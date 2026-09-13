package computer

// Commands is the set of root command groups.
type Computer struct {
	Read Read `command:"read" alias:"r" description:"Read a Computer object in Landscape."`
	// Create   Create `command:"create" alias:"c" description:"Create a new MidPoint user object."`
	// Search   Search `command:"search" alias:"s" description:"Search for MidPoint user objects."`
	// Delete   Delete `command:"delete" alias:"d" description:"Delete a MidPoint user object."`
	// Password struct {
	// 	Policy Policy `command:"policy" alias:"p" description:"Show a MidPoint user object password policy."`
	// 	Reset  Reset  `command:"reset" alias:"r" description:"Reset a MidPoint user password."`
	// } `command:"password" alias:"p" description:"Operate on MidPoint users' passwords."`
}
