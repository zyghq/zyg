package zyg

import (
	"fmt"
	"os"
	"strconv"
)

const DefaultSecretKeyLength = 64

func GetEnv(key string) (string, error) {
	value, status := os.LookupEnv(key)
	if !status {
		return "", fmt.Errorf("env `%s` is not set", key)
	}
	return value, nil
}

func GetAvatarBaseURL() string {
	value, ok := os.LookupEnv("AVATAR_BASE_URL")
	if !ok {
		return "https://avatar.vercel.sh/" // probably self-host?
	}
	return value
}

func DBQueryDebug() bool {
	debug, err := strconv.ParseBool(os.Getenv("ZYG_DB_QUERY_DEBUG"))
	if err != nil {
		return false
	}
	return debug
}
