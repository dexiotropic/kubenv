package version

import "fmt"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// String returns the version string with commit and build date metadata.
func String() string {
	return fmt.Sprintf("%s (commit %s, built %s)", version, commit, date)
}

// Commit returns the build commit.
func Commit() string {
	return commit
}

// Date returns the build date.
func Date() string {
	return date
}
