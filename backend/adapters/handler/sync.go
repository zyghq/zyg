package handler

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SyncHandler struct {
	ws  ports.WorkspaceServicer
	ths ports.ThreadServicer
}

func NewSyncHandler(ws ports.WorkspaceServicer, ths ports.ThreadServicer) *SyncHandler {
	return &SyncHandler{
		ws:  ws,
		ths: ths,
	}
}

func (h *SyncHandler) syncWorkspaceMemberShapesV1(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()
	hub := sentry.GetHubFromContext(ctx)

	// Upstream URL is taken directly from the Electric docs
	// see: https://electric-sql.com/docs/guides/auth
	// configured for base URL.
	reqURL := r.URL
	originURL, err := url.Parse(zyg.ElectricBaseUrl() + "/v1/shape")
	if err != nil {
		hub.CaptureException(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Copy allowed query params
	allowedParams := map[string]bool{
		"live": true, "table": true, "handle": true,
		"offset": true, "cursor": true,
	}

	// Create a new url.Values to store query parameters
	q := originURL.Query()
	q.Set("where", fmt.Sprintf("workspace_id = '%s'", member.WorkspaceId))

	// Copy the allowed parameters from the request
	for key, values := range reqURL.Query() {
		if allowedParams[key] {
			q.Set(key, values[0])
		}
	}
	originURL.RawQuery = q.Encode()

	// Make upstream request
	client := &http.Client{}
	upstreamReq, err := http.NewRequest(http.MethodGet, originURL.String(), nil)
	if err != nil {
		hub.CaptureException(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(upstreamReq)
	if err != nil {
		hub.CaptureException(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// Copy response headers except content-encoding and content-length
	for key, values := range resp.Header {
		if !strings.EqualFold(key, "Content-Encoding") &&
			!strings.EqualFold(key, "Content-Length") {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// Copy status code and body
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		hub.CaptureException(err)
		// Log error but response is already started
		fmt.Printf("Error copying response: %v\n", err)
	}
}

func (h *SyncHandler) syncWorkspaceCustomerShapesV1(
	w http.ResponseWriter, r *http.Request, member *models.Member) {
	ctx := r.Context()
	hub := sentry.GetHubFromContext(ctx)

	// Upstream URL is taken directly from the Electric docs
	// see: https://electric-sql.com/docs/guides/auth
	// configured for base URL.
	reqURL := r.URL
	originURL, err := url.Parse(zyg.ElectricBaseUrl() + "/v1/shape")
	if err != nil {
		hub.CaptureException(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Copy allowed query params
	allowedParams := map[string]bool{
		"live": true, "table": true, "handle": true,
		"offset": true, "cursor": true,
	}

	// Create a new url.Values to store query parameters
	q := originURL.Query()
	q.Set("where", fmt.Sprintf("workspace_id = '%s'", member.WorkspaceId))

	// Copy the allowed parameters from the request
	for key, values := range reqURL.Query() {
		if allowedParams[key] {
			q.Set(key, values[0])
		}
	}
	originURL.RawQuery = q.Encode()

	// Make upstream request
	client := &http.Client{}
	upstreamReq, err := http.NewRequest(http.MethodGet, originURL.String(), nil)
	if err != nil {
		hub.CaptureException(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(upstreamReq)
	if err != nil {
		hub.CaptureException(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// Copy response headers except content-encoding and content-length
	for key, values := range resp.Header {
		if !strings.EqualFold(key, "Content-Encoding") &&
			!strings.EqualFold(key, "Content-Length") {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// Copy status code and body
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		hub.CaptureException(err)
		// Log error but response is already started
		fmt.Printf("Error copying response: %v\n", err)
	}
}
