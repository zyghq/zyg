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

func GetXServerUrl() string {
	value, ok := os.LookupEnv("ZYG_XSERVER_URL")
	if !ok {
		return "http://localhost:8000"
	}
	return value
}

func LandingPageUrl() string {
	value, ok := os.LookupEnv("ZYG_URL")
	if !ok {
		return "https://zyg.ai"
	}
	return value
}

func ResendApiKey() string {
	value, ok := os.LookupEnv("RESEND_API_KEY")
	if !ok {
		return ""
	}
	return value
}

func CFAccountId() string {
	value, ok := os.LookupEnv("CF_ACCOUNT_ID")
	if !ok {
		return ""
	}
	return value
}

func R2AccessKeyId() string {
	value, ok := os.LookupEnv("R2_ACCESS_KEY_ID")
	if !ok {
		return ""
	}
	return value
}

func R2AccessSecretKey() string {
	value, ok := os.LookupEnv("R2_ACCESS_SECRET_KEY")
	if !ok {
		return ""
	}
	return value
}

func S3Bucket() string {
	value, ok := os.LookupEnv("S3_BUCKET")
	if !ok {
		return "zygdev"
	}
	return value
}
