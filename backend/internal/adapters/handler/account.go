package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/zyghq/zyg/internal/domain"
	"github.com/zyghq/zyg/internal/ports"
)

type AccountHandler struct {
	as ports.AccountServicer
}

func NewAccountHandler(as ports.AccountServicer) *AccountHandler {
	return &AccountHandler{as: as}
}

func (h *AccountHandler) handleGetPatList(w http.ResponseWriter, r *http.Request, account *domain.Account) {
	ctx := r.Context()
	aps, err := h.as.GetUserPatList(ctx, account.AccountId)
	if err != nil {
		slog.Error("failed to get list of pats "+
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
				"might need to check the json encoding defn",
			slog.String("accountId", account.AccountId),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
