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

type AuthenticatedAccountHandler func(http.ResponseWriter, *http.Request, *models.Account)

type AuthenticatedMemberHandler func(http.ResponseWriter, *http.Request, *models.Member)

//var (
//	key   = []byte(zyg.WorkOSCookiePassword())
//	sessionStore = sessions.NewCookieStore(key)
//)

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

		account, err := authz.AuthenticateUserAccount(ctx, sub)
		if errors.Is(err, services.ErrAccountNotFound) {
			slog.Error("auth account not found", slog.Any("error", err))
			return account, fmt.Errorf("account not found or does not exist")
		}
		if errors.Is(err, services.ErrAccount) {
			slog.Error("failed to fetch account by auth user id", slog.Any("error", err))
			return account, fmt.Errorf("failed to get account by auth user id: %s got error: %v", sub, err)
		}
		if err != nil {
			slog.Error("something went wrong", slog.Any("error", err))
			return account, fmt.Errorf("failed to get account by auth user id: %s got error: %v", sub, err)
		}
		return account, nil
	} else {
		return models.Account{}, fmt.Errorf("unsupported scheme `%s` cannot authenticate", scheme)
	}
}

func AuthenticateMember(
	ctx context.Context, authz ports.AuthServicer, workspaceId string, scheme string, cred string,
) (models.Member, error) {
	if scheme == "bearer" {
		account, err := AuthenticateAccount(ctx, authz, scheme, cred)
		if err != nil {
			return models.Member{}, err
		}

		member, err := authz.AuthenticateWorkspaceMember(ctx, workspaceId, account.AccountId)
		if err != nil {
			return models.Member{}, fmt.Errorf("failed to authenticate workspace member: %v", err)
		}

		return member, nil
	} else {
		return models.Member{}, fmt.Errorf("unsupported scheme `%s` cannot authenticate", scheme)
	}
}

//func handleWorkOSAuthLogin(w http.ResponseWriter, r *http.Request) {
//	apiKey := zyg.WorkOSAPIKey()
//	clientID := zyg.WorkOSClientID()
//
//	usermanagement.SetAPIKey(apiKey)
//	url, err := usermanagement.GetAuthorizationURL(usermanagement.GetAuthorizationURLOpts{
//		Provider:    "authkit",
//		ClientID:    clientID,
//		RedirectURI: fmt.Sprintf("%s/auth/callback/", zyg.SrvBaseURL()),
//	})
//	if err != nil {
//		slog.Error("failed to get auth url from workos", slog.Any("error", err))
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//		return
//	}
//	http.Redirect(w, r, url.String(), http.StatusSeeOther)
//}

//func handleWorkOSAuthCallback(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//	apiKey := zyg.WorkOSAPIKey()
//	clientID := zyg.WorkOSClientID()
//
//	usermanagement.SetAPIKey(apiKey)
//	code := r.URL.Query().Get("code")
//	if code == "" {
//		slog.Error("code is empty")
//		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//		return
//	}
//
//	response, err := usermanagement.AuthenticateWithCode(
//		ctx,
//		usermanagement.AuthenticateWithCodeOpts{
//			ClientID: clientID,
//			Code:     code,
//		},
//	)
//	if err != nil {
//		slog.Error("failed to authenticate with code", slog.Any("error", err))
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//		return
//	}
//
//	session, _ := sessionStore.Get(r, "wos-session")
//	session.Options = &sessions.Options{
//		Path:     "/",
//		HttpOnly: true,
//		Secure:   false,
//		SameSite: http.SameSiteLaxMode,
//	}
//	session.Values["authenticated"] = true
//	if err := session.Save(r, w); err != nil {
//		slog.Error("failed to save cookie session", slog.Any("error", err))
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//		return
//	}
//
//	fmt.Println("******************** auth profile ***********************")
//	fmt.Println(response)
//	fmt.Println("******************** auth profile ***********************")
//
//	// redirect to frontend app URL
//	redirectURL := fmt.Sprintf("%s/workspaces/", zyg.AppBaseURL())
//	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
//}
