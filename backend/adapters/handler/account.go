package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
	"github.com/zyghq/zyg/services"
)

type AccountHandler struct {
	as ports.AccountServicer
	ws ports.WorkspaceServicer
}

func NewAccountHandler(as ports.AccountServicer, ws ports.WorkspaceServicer) *AccountHandler {
	return &AccountHandler{as: as, ws: ws}
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
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else if scheme == "bearer" {
		hmacSecret, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			slog.Error("env SUPABASE_JWT_SECRET is not set")
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
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var reqp map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&reqp)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var accountName string
		if name, found := reqp["name"]; found {
			if name == nil {
				slog.Error(
					"name cannot be empty", slog.String("subject", sub),
				)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			ns := name.(string)
			accountName = ns
		}

		account, isCreated, err := h.as.CreateAuthAccount(ctx, sub, ac.Email, accountName, services.DefaultAuthProvider)
		if err != nil {
			slog.Error(
				"failed to create auth account", slog.Any("err", err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// @sanchitrk
		// add CDP event when the account is created
		if isCreated {
			slog.Info("created auth account", slog.String("accountId", account.AccountId))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(account); err != nil {
				slog.Error("failed to encode json", slog.Any("err", err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(account); err != nil {
				slog.Error("failed to encode json", slog.Any("err", err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	} else {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
}

func (h *AccountHandler) handleGetPatList(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()
	aps, err := h.as.ListPersonalAccessTokens(ctx, account.AccountId)
	if err != nil {
		slog.Error("failed to fetch pat list", slog.Any("err", err))
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

func (h *AccountHandler) handleCreatePat(w http.ResponseWriter, r *http.Request, account *models.Account) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()

	var reqp PATReq
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pat, err := h.as.GeneratePersonalAccessToken(ctx, account.AccountId, reqp.Name, reqp.Description)
	if err != nil {
		slog.Error(
			"failed to create pat", slog.Any("err", err),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pat); err != nil {
		slog.Error("failed to encode json", slog.String("patId", pat.PatId))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *AccountHandler) handleDeletePat(w http.ResponseWriter, r *http.Request, account *models.Account) {
	ctx := r.Context()
	patId := r.PathValue("patId")

	pat, err := h.as.GetPersonalAccessToken(ctx, patId)
	if errors.Is(err, repository.ErrEmpty) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to fetch pat for deletion", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.as.DeletePersonalAccessToken(ctx, pat.PatId)
	if err != nil {
		slog.Error("failed to delete pat", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
