package scanner

// All is the master list of every scanner.
// Adding a new scanner = add one line here. No other changes needed.
var All = []Scanner{
	&PasswordSSH{},
	&AccessibleRDP{},
	&AccessibleDB{},
}
