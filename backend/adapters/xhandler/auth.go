package xhandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
	"github.com/zyghq/zyg/services"
)

func CheckAuthCredentials(r *http.Request) (string, string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", "", fmt.Errorf("no authorization header provided")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid token")
	}
	scheme := strings.ToLower(parts[0])
	return scheme, parts[1], nil
}

func AuthenticateCustomer(
	ctx context.Context, authz ports.CustomerAuthServicer,
	scheme string, cred string,
) (models.Customer, error) {
	var customer models.Customer
	if scheme == "bearer" {
		slog.Info("authenticate with customer JWT")
		hmacSecret, err := zyg.GetEnv("ZYG_CUSTOMER_JWT_SECRET")
		if err != nil {
			return customer, fmt.Errorf("failed to get env SUPABASE_JWT_SECRET with error: %v", err)
		}

		cc, err := services.ParseCustomerJWTToken(cred, []byte(hmacSecret))
		if err != nil {
			return customer, fmt.Errorf("failed to parse JWT token with error: %v", err)
		}

		sub, err := cc.RegisteredClaims.GetSubject()
		if err != nil {
			return customer, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}

		slog.Info("authenticated customer with customer id", slog.String("customerId", sub))

		customer, err = authz.WorkspaceCustomer(ctx, cc.WorkspaceId, sub)

		if errors.Is(err, services.ErrCustomerNotFound) {
			slog.Warn(
				"customer not found or does not exist",
				slog.String("customerId", sub),
			)
			return customer, fmt.Errorf("customer not found or does not exist")
		}

		if err != nil {
			slog.Error(
				"failed to get customer by customer id"+
					"something went wrong",
				slog.String("customerId", sub),
			)
			return customer, fmt.Errorf("failed to get customer by customer id: %s got error: %v", sub, err)
		}

		return customer, nil
	} else {
		return customer, fmt.Errorf("unsupported scheme: `%s` cannot authenticate", scheme)
	}
}

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *models.Customer)
