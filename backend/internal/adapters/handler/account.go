package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/internal/adapters/repository"
	"github.com/zyghq/zyg/internal/domain"
	"github.com/zyghq/zyg/internal/ports"
	"github.com/zyghq/zyg/internal/services"
)

type AccountHandler struct {
	as ports.AccountServicer
}

func NewAccountHandler(as ports.AccountServicer) *AccountHandler {
	return &AccountHandler{as: as}
}

func (h *AccountHandler) handleGetOrCreateAccount(w http.ResponseWriter, r *http.Request) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	scheme, cred, err := CheckAuthCredentials(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	if scheme == "token" {
		slog.Warn("token authorization scheme unsupported for auth account creation")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else if scheme == "bearer" {
		hmacSecret, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			slog.Error(
				"failed to get env SUPABASE_JWT_SECRET " +
					"required to decode the incoming jwt token",
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		ac, err := services.ParseJWTToken(cred, []byte(hmacSecret))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		sub, err := ac.RegisteredClaims.GetSubject()
		if err != nil {
			slog.Warn("failed to get subject from parsed token - make sure it is set in the token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// initialize auth user account by subject.
		account := domain.Account{AuthUserId: sub, Email: ac.Email, Provider: services.DefaultAuthProvider}
		account, isCreated, err := h.as.InitiateAccount(ctx, account)
		if err != nil {
			slog.Error(
				"failed to get or create account by subject "+
					"perhaps a failed query or mapping", slog.String("subject", sub),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if isCreated {
			slog.Info("created auth account for subject", slog.String("accountId", account.AccountId))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(account); err != nil {
				slog.Error(
					"failed to encode account to json "+
						"check the json encoding defn",
					slog.String("accountId", account.AccountId),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		} else {
			slog.Info("auth account already exists for subject", slog.String("accountId", account.AccountId))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(account); err != nil {
				slog.Error(
					"failed to encode account to json "+
						"check the json encoding defn",
					slog.String("accountId", account.AccountId),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	} else {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
}

func (h *AccountHandler) handleGetPatList(w http.ResponseWriter, r *http.Request, account *domain.Account) {
	ctx := r.Context()
	aps, err := h.as.UserPats(ctx, account.AccountId)
	if err != nil {
		slog.Error("failed to pat list "+
			"something went wrong",
			slog.String("accountId", account.AccountId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(aps); err != nil {
		slog.Error(
			"failed to encode pats to json "+
				"check the json encoding defn",
			slog.String("accountId", account.AccountId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *AccountHandler) handleCreatePat(w http.ResponseWriter, r *http.Request, account *domain.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	var rb PATReqPayload
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ap := domain.AccountPAT{
		AccountId:   account.AccountId,
		Name:        rb.Name,
		UnMask:      true, // unmask only once created
		Description: domain.NullString(rb.Description),
	}

	ap, err = h.as.IssuePersonalAccessToken(ctx, ap)
	if err != nil {
		slog.Error(
			"failed to create account PAT "+
				"perhaps check the db query or mapping",
			slog.String("accountId", account.AccountId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ap); err != nil {
		slog.Error(
			"failed to encode account pat to json "+
				"check the json encoding defn",
			slog.String("patId", ap.PatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *AccountHandler) handleDeletePat(w http.ResponseWriter, r *http.Request, account *domain.Account) {
	ctx := r.Context()
	patId := r.PathValue("patId")

	pat, err := h.as.UserPat(ctx, patId)
	if errors.Is(err, repository.ErrEmpty) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to get pat by pat id "+
			"something went wrong",
			slog.String("patId", patId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.as.HardDeletePat(ctx, pat.PatId)
	if err != nil {
		slog.Error("failed to delete pat "+
			"something went wrong",
			slog.String("patId", pat.PatId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// return http status 204 no content
	w.WriteHeader(http.StatusNoContent)
}
