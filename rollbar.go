package code

import (
	"fmt"
	"os"

	"github.com/rollbar/rollbar-go"
)

var commitHash string

func setupRollbar() {
	token := getEnv("ROLLBAR_TOKEN", "")

	if token == "" {
		return
	}

	fmt.Println("ROLLBAR_TOKEN: ", token)
	fmt.Println("ROLLBAR_ENV: ", getEnv("ROLLBAR_ENV", "development"))
	fmt.Println("commitHash: ", commitHash)
	fmt.Println("HOSTNAME: ", getEnv("HOSTNAME", ""))
	fmt.Println("ROLLBAR_SERVER_ROOT: ", getEnv("ROLLBAR_SERVER_ROOT", "/"))

	rollbar.SetToken(token)
	rollbar.SetEnvironment(getEnv("ROLLBAR_ENV", "development"))
	rollbar.SetCodeVersion(commitHash)
	rollbar.SetServerHost(getEnv("HOSTNAME", ""))
	rollbar.SetServerRoot(getEnv("ROLLBAR_SERVER_ROOT", "/"))
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
