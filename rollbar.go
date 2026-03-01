package main

import (
	"os"

	"github.com/rollbar/rollbar-go"
)

var commitHash string

func setupRollbar() bool {
	token := getEnv("ROLLBAR_TOKEN", "")

	if token == "" {
		return false
	}

	rollbar.SetToken(token)
	rollbar.SetEnvironment(getEnv("ROLLBAR_ENV", "development"))
	rollbar.SetCodeVersion(commitHash)
	rollbar.SetServerHost(getEnv("HOSTNAME", ""))
	rollbar.SetServerRoot(getEnv("ROLLBAR_SERVER_ROOT", "/"))

	return true
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
