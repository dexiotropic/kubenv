package version

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func String() string {
	return version
}

func Commit() string {
	return commit
}

func Date() string {
	return date
}
