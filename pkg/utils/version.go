package utils

import "fmt"

var (
	version        = "unknown" // Filled out during build
	commitHash     = "0000000"
	buildTimestamp = "0000-00-00T00:00:00"
	dirty          = ""
)

func GetShortVersion() string {
	return version
}

func GetLongVersion() string {
	verStr := fmt.Sprintf("%s_%s", version, commitHash)
	if dirty != "" {
		verStr += "-dirty"
	}
	return verStr
}

func GetBuildTimestamp() string {
	return buildTimestamp
}
