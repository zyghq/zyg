package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type UserHandler struct {
	us ports.UserServicer
}

func NewUserHandler(userService ports.UserServicer) *UserHandler {
	return &UserHandler{
		us: userService,
	}
}

func (h *UserHandler) handleWorkOSWebhook(w http.ResponseWriter, r *http.Request) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	ctx := r.Context()
	hub := sentry.GetHubFromContext(ctx)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to read body", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	workOSEvent, err := models.InitWorkOSEventFromPayload(payload)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to initialize workos user", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fmt.Println("WORKOS EVENT", workOSEvent.Event)

	user, err := h.us.CreateWorkOSUser(ctx, &workOSEvent.WorkOSUser)
	if err != nil {
		hub.CaptureException(err)
		slog.Error("failed to create workos user", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		slog.Error("failed to encode json", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
