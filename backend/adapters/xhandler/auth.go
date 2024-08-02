package xhandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

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
	scheme string, cred string, widgetId string) (models.Customer, error) {
	var customer models.Customer
	if scheme == "bearer" {

		sk, err := authz.GetWidgetLinkedSecretKey(ctx, widgetId)
		if err != nil {
			return customer, fmt.Errorf("%v", err)
		}

		cc, err := services.ParseCustomerJWTToken(cred, []byte(sk.SecretKey))
		if err != nil {
			return customer, fmt.Errorf("%v", err)
		}

		sub, err := cc.RegisteredClaims.GetSubject()
		if err != nil {
			return customer, fmt.Errorf("%v", err)
		}

		customer, err = authz.GetWorkspaceCustomerIgnoreRole(ctx, cc.WorkspaceId, sub)
		if errors.Is(err, services.ErrCustomerNotFound) {
			return customer, fmt.Errorf("customer not found or does not exist")
		}

		if err != nil {
			slog.Error("failed to fetch customer", slog.Any("err", err))
			return customer, fmt.Errorf("failed to validate customer with customer id: %s got error: %v", sub, err)
		}

		return customer, nil
	} else {
		return customer, fmt.Errorf("unsupported scheme: `%s` cannot authenticate", scheme)
	}
}

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *models.Customer)
