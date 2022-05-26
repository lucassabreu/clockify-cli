package cmdutil

// Version register which is the CLI tag, commit and build date
type Version struct {
	Tag    string
	Commit string
	Date   string
}
