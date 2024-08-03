package handler

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

func AuthenticateAccount(
	ctx context.Context, authz ports.AuthServicer, scheme string, cred string) (models.Account, error) {
	if scheme == "token" {
		account, err := authz.ValidatePersonalAccessToken(ctx, cred)
		if err != nil {
			return account, fmt.Errorf("failed to authenticate got error: %v", err)
		}
		slog.Info("authenticated account with PAT", slog.String("accountId", account.AccountId))
		return account, nil
	} else if scheme == "bearer" {
		hmacSecret, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			return models.Account{}, fmt.Errorf("failed to get env SUPABASE_JWT_SECRET got error: %v", err)
		}
		ac, err := services.ParseJWTToken(cred, []byte(hmacSecret))
		if err != nil {
			return models.Account{}, fmt.Errorf("failed to parse JWT token got error: %v", err)
		}
		sub, err := ac.RegisteredClaims.GetSubject()
		if err != nil {
			return models.Account{}, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}

		slog.Info("authenticated account with jwt", slog.String("authUserId", sub))
		account, err := authz.AuthenticateUser(ctx, sub)

		if errors.Is(err, services.ErrAccountNotFound) {
			return account, fmt.Errorf("account not found or does not exist")
		}
		if errors.Is(err, services.ErrAccount) {
			slog.Error(
				"failed to get account by auth user id "+
					"perhaps a failed query or mapping",
				slog.String("authUserId", sub),
			)
			return account, fmt.Errorf("failed to get account by auth user id: %s got error: %v", sub, err)
		}
		if err != nil {
			slog.Error(
				"failed to get account by auth user id "+
					"something went wrong",
				slog.String("authUserId", sub),
			)
			return account, fmt.Errorf("failed to get account by auth user id: %s got error: %v", sub, err)
		}
		return account, nil
	} else {
		return models.Account{}, fmt.Errorf("unsupported scheme `%s` cannot authenticate", scheme)
	}
}

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *models.Account)
