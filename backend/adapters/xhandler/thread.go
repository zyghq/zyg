package xhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/services"
	"io"
	"log/slog"
	"net/http"
)

func (h *ThreadHandler) handlePostmarkInboundMessage(w http.ResponseWriter, r *http.Request) {
	defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	var reqp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&reqp)
	if err != nil {
		slog.Error("error decoding json payload", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	workspaceId := r.PathValue("workspaceId")
	workspace, err := h.ws.GetWorkspace(ctx, workspaceId)
	if errors.Is(err, services.ErrWorkspaceNotFound) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if err != nil {
		slog.Error("failed to fetch workspace", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Info("got postmark inbound message for workspace",
		slog.Any("workspaceId", workspace.WorkspaceId),
	)

	inboundMessage, err := (&models.PostmarkInboundMessage{}).FromPayload(reqp)
	if err != nil {
		slog.Error("error parsing postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println("*********** Postmark Inbound Message ***********")
	fmt.Println(inboundMessage.Subject())
	fmt.Println(inboundMessage.PlainText())
	fmt.Println(inboundMessage.Html())
	fmt.Println("*********************************************")

	// Pull Customer mail and name from the inbound mail and get or create the Customer.
	fromEmail := inboundMessage.FromEmail()
	fromName := inboundMessage.FromName()
	customer, _, err := h.ws.CreateCustomerWithEmail(
		ctx, workspace.WorkspaceId, fromEmail, true, fromName)
	if err != nil {
		slog.Error("error creating customer for postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Get the system Member from the Workspace which will process this inbound mail.
	member, err := h.ws.GetSystemMember(ctx, workspace.WorkspaceId)
	if err != nil {
		slog.Error("error getting system member for postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// process the postmark inbound message.
	thread, message, err := h.ths.ProcessPostmarkInbound(
		ctx, workspace.WorkspaceId, customer.AsCustomerActor(),
		member.AsMemberActor(), inboundMessage,
	)
	if err != nil {
		slog.Error("error processing postmark inbound message", slog.Any("err", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println("********* thread + message **********")
	fmt.Println("ThreadID", thread.ThreadId)
	fmt.Println("MessageID", message.MessageId)
	fmt.Println("Thread Title", thread.Title)
	fmt.Println("*********************************************")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("ok"))
	if err != nil {
		return
	}
}
