package log

import (
	"os"

	logging "github.com/ipfs/go-log/v2"
)

// nolint
func SetupLogLevels() {
	if _, set := os.LookupEnv("GOLOG_LOG_LEVEL"); !set {
		logging.SetLogLevel("*", "INFO")
	}
}
