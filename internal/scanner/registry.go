package scanner

// All is the master list of every scanner.
// Adding a new scanner requires only a single line here.
var All = []Scanner{
	&PasswordSSH{},
	&AccessibleRDP{},
	&AccessibleDB{},
}

// ByName maps each scanner name to its implementation for O(1) lookup.
var ByName = func() map[string]Scanner {
	m := make(map[string]Scanner, len(All))
	for _, s := range All {
		m[s.Name()] = s
	}
	return m
}()

// Defaults returns all scanners that are enabled by default (DefaultEnabled() == true).
func Defaults() []Scanner {
	var out []Scanner
	for _, s := range All {
		if s.DefaultEnabled() {
			out = append(out, s)
		}
	}
	return out
}
