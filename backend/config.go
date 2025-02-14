package zyg

import (
	"fmt"
	"github.com/google/uuid"
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

func ServerDomain() string {
	value, ok := os.LookupEnv("ZYG_SERVER_DOMAIN")
	if !ok {
		return "localhost"
	}
	return value
}

func ServerProto() string {
	value, ok := os.LookupEnv("ZYG_SERVER_PROTO")
	if !ok {
		return "http"
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

func RedisAddr() string {
	value, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		return "localhost:6379"
	}
	return value
}

func RedisUsername() string {
	value, ok := os.LookupEnv("REDIS_USERNAME")
	if !ok {
		return "zygdev"
	}
	return value
}

func RedisPassword() string {
	value, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		return ""
	}
	return value
}

func RedisTLSEnabled() bool {
	enabled, err := strconv.ParseBool(os.Getenv("REDIS_TLS_ENABLED"))
	if err != nil {
		return false
	}
	return enabled
}

func SentryDebugEnabled() bool {
	enabled, err := strconv.ParseBool(os.Getenv("SENTRY_DEBUG_ENABLED"))
	if err != nil {
		return false
	}
	return enabled
}

func SentryEnv() string {
	value, ok := os.LookupEnv("SENTRY_ENV")
	if !ok {
		return "staging"
	}
	return value
}

func PostmarkAccountToken() string {
	value, ok := os.LookupEnv("POSTMARK_ACCOUNT_TOKEN")
	if !ok {
		return ""
	}
	return value
}

// WebhookUsername retrieves the "WEBHOOK_USERNAME" environment variable or generates a UUID if not found.
// If not set then generate random username for security reasons.
func WebhookUsername() string {
	value, ok := os.LookupEnv("WEBHOOK_USERNAME")
	if !ok {
		u, _ := uuid.NewUUID()
		return u.String()
	}
	return value
}

// WebhookPassword retrieves the "WEBHOOK_PASSWORD" environment variable or generates a new UUID if not set.
// If not set then generate random username for security reasons.
func WebhookPassword() string {
	value, ok := os.LookupEnv("WEBHOOK_PASSWORD")
	if !ok {
		u, _ := uuid.NewUUID()
		return u.String()
	}
	return value
}

func WebhookUrl() string {
	proto := ServerProto()
	domain := ServerDomain()
	u := WebhookUsername()
	p := WebhookPassword()
	return fmt.Sprintf("%s://%s:%s@%s", proto, u, p, domain)
}

func PostmarkDeliveryDomain() string {
	value, ok := os.LookupEnv("POSTMARK_DELIVERY_DOMAIN")
	if !ok {
		return "mtasv.net"
	}
	return value
}
