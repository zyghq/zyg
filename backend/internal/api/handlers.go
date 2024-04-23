package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/internal/auth"
	"github.com/zyghq/zyg/internal/model"
)

func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// request
type LLMRREval struct {
	Eval int `json:"eval"`
}

// response
type LLMResponse struct {
	Text      string `json:"text"`
	RequestId string `json:"requestId"`
	Model     string `json:"model"`
}

type LLM struct {
	WorkspaceId string
	Prompt      string
	RequestId   string
}

// for now we are directly using the Ollama server
// later will update to use our `converse` server probably with grpc.
// similarly the `LLMResponse` will be updated to include other specific fields.
func (llm LLM) Generate() (LLMResponse, error) {

	var err error
	var response LLMResponse

	buf := new(bytes.Buffer)
	// for now this is specific to the Ollama server
	// will update to use our `converse` server probably with grpc.
	body := struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}{
		Model:  "llama2",
		Prompt: llm.Prompt,
		Stream: false,
	}

	err = json.NewEncoder(buf).Encode(&body)
	if err != nil {
		log.Printf("failed to encode LLM request body for requestId: %s with error: %v", llm.RequestId, err)
		return response, err
	}

	log.Printf("LLM request for workspaceId: %s with requestId: %s", llm.WorkspaceId, llm.RequestId)
	resp, err := http.Post("http://0.0.0.0:11434/api/generate", "application/json", buf)
	if err != nil {
		log.Printf("LLM request failed for requestId: %s with error: %v", llm.RequestId, err)
		return response, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("expected status %d; but got %d", http.StatusOK, resp.StatusCode)
	}

	// response structure from Ollama server
	// will be updated based on the `converse` server response.
	rb := struct {
		Model              string `json:"model"`
		CreatedAt          string `json:"created_at"`
		Response           string `json:"response"`
		Done               bool   `json:"done"`
		Context            []int  `json:"context"`
		TotalDuration      int    `json:"total_duration"`
		LoadDuration       int    `json:"load_duration"`
		PromptEvalCount    int    `json:"prompt_eval_count"`
		PromptEvalDuration int    `json:"prompt_eval_duration"`
		EvalCount          int    `json:"eval_count"`
		EvalDuration       int    `json:"eval_duration"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&rb)
	if err != nil {
		log.Printf("failed to decode LLM response for requestId: %s with error: %v\n", llm.RequestId, err)
		return response, err
	}

	return LLMResponse{
		Text:      rb.Response,
		RequestId: llm.RequestId,
		Model:     rb.Model,
	}, err
}

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleGetOrCreateAuthAccount(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		scheme, cred, err := HttpAuthCredentials(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

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
			ac, err := auth.ParseJWTToken(cred, []byte(hmacSecret))
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
			// get or create account by subject maps which auth user id
			account := model.Account{AuthUserId: sub, Email: ac.Email, Provider: auth.DefaultAuthProvider}
			account, isCreated, err := account.GetOrCreateByAuthUserId(ctx, db)
			if errors.Is(err, model.ErrQuery) || errors.Is(err, model.ErrEmpty) {
				slog.Error(
					"failed to get or create account by subject "+
						"perhaps a failed query or mapping", slog.String("subject", sub),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if err != nil {
				slog.Error("failed to get or create account by subject "+
					"something went wrong", slog.String("subject", sub),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if isCreated {
				slog.Info("successfully created auth account for subject", slog.String("accountId", account.AccountId))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				if err := json.NewEncoder(w).Encode(account); err != nil {
					slog.Error(
						"failed to encode account to json "+
							"might need to check the json encoding defn",
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
							"might need to check the json encoding defn",
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
	})
}

func handleCreatePAT(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var rb PATReqPayload
		err := json.NewDecoder(r.Body).Decode(&rb)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		account, err := AuthenticateAccount(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ap := model.AccountPAT{
			AccountId:   account.AccountId,
			Name:        rb.Name,
			UnMask:      true, // unmask only once created
			Description: NullString(rb.Description),
		}

		ap, err = ap.Create(ctx, db)
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
					"might need to check the json encoding defn",
				slog.String("patId", ap.PatId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetPATs(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, err := AuthenticateAccount(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ap := model.AccountPAT{AccountId: account.AccountId}
		aps, err := ap.GetListByAccountId(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"no account PATs found for account",
				slog.String("accountId", account.AccountId),
			)
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
			return
		}
		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get account PATs "+
					"perhaps a failed query or mapping",
				slog.String("accountId", account.AccountId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
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
	})
}

func handleCreateWorkspace(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var rb WorkspaceReqPayload
		err := json.NewDecoder(r.Body).Decode(&rb)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		account, err := AuthenticateAccount(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		workspace := model.Workspace{AccountId: account.AccountId, Name: rb.Name}
		workspace, err = workspace.Create(ctx, db)

		// checks if there was a db query error
		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to create workspace "+
					"perhaps check the db query",
				slog.String("accountId", account.AccountId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// handle all errors
		if err != nil {
			slog.Error(
				"failed to create workspace "+
					"something went wrong",
				slog.String("accountId", account.AccountId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// everything went well return the newly created workspace
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(workspace); err != nil {
			slog.Error(
				"failed to encode workspace to json "+
					"might need to check the json encoding defn",
				slog.String("workspaceId", workspace.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetWorkspaces(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		account, err := AuthenticateAccount(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspace := model.Workspace{AccountId: account.AccountId}

		// get list of workspaces for account
		workspaces, err := workspace.GetListByAccountId(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"no workspaces found for account",
				slog.String("accountId", account.AccountId),
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(workspaces); err != nil {
				slog.Error(
					"failed to encode workspaces to json "+
						"might need to check the json encoding defn",
					slog.String("accountId", account.AccountId),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get workspaces "+
					"perhaps a failed query or mapping",
				slog.String("accountId", account.AccountId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get list of workspaces "+
					"something went wrong",
				slog.String("accountId", account.AccountId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(workspaces); err != nil {
			slog.Error(
				"failed to encode workspaces to json "+
					"might need to check the json encoding defn",
				slog.String("accountId", account.AccountId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetWorkspace(ctx context.Context, db *pgxpool.Pool) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		account, err := AuthenticateAccount(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspaceId := r.PathValue("workspaceId")

		workspace := model.Workspace{WorkspaceId: workspaceId, AccountId: account.AccountId}
		workspace, err = workspace.GetAccountWorkspace(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"account workspace not found or does not exist",
				slog.String("workspaceId", workspaceId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get account workspace or does not exist "+
					"perhaps check the db query or mapping",
				slog.String("workspaceId", workspaceId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get account workspace or does not exist "+
					"something went wrong",
				slog.String("workspaceId", workspaceId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(workspace); err != nil {
			slog.Error(
				"failed to encode workspace to json "+
					"might need to check the json encoding defn",
				slog.String("workspaceId", workspace.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetCustomer(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customer, err := AuthenticateCustomer(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(customer); err != nil {
			slog.Error(
				"failed to encode customer to json "+
					"might need to check the json encoding defn",
				slog.String("customerId", customer.CustomerId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleCustomerTokenIssue(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var rb CustomerTIReqPayload
		err := json.NewDecoder(r.Body).Decode(&rb)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		externalId := NullString(rb.Customer.ExternalId)
		email := NullString(rb.Customer.Email)
		phone := NullString(rb.Customer.Phone)
		if !externalId.Valid && !email.Valid && !phone.Valid {
			slog.Error("at least one of `externalId`, `email` or `phone` is required")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// authenticate account
		account, err := AuthenticateAccount(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspaceId := r.PathValue("workspaceId")

		workspace := model.Workspace{WorkspaceId: workspaceId, AccountId: account.AccountId}
		workspace, err = workspace.GetAccountWorkspace(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"account workspace not found or does not exist",
				slog.String("workspaceId", workspaceId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get account workspace or does not exist "+
					"perhaps check the db query or mapping",
				slog.String("workspaceId", workspaceId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get account workspace or does not exist "+
					"something went wrong",
				slog.String("workspaceId", workspaceId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		customer := model.Customer{
			WorkspaceId: workspace.WorkspaceId,
			ExternalId:  externalId,
			Email:       email,
			Phone:       phone,
		}

		var isCreated bool
		var resp CustomerTIRespPayload

		if rb.Create {
			if rb.CreateBy == nil {
				slog.Error("requires `createBy` when `create` is enabled")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			createBy := *rb.CreateBy
			slog.Info("create Customer if does not exists", slog.String("createBy", createBy))
			switch createBy {
			case "email":
				if !customer.Email.Valid {
					slog.Error("email is required for createBy email")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}

				slog.Info("create Customer by email")
				customer, isCreated, err = customer.GetOrCreateWrkCustomerByEmail(ctx, db)

				if errors.Is(err, model.ErrEmpty) {
					slog.Warn(
						"Customer not found or does not exist after creating" +
							"perhaps the customer is not created or is not returned",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				if errors.Is(err, model.ErrQuery) {
					slog.Error(
						"failed to get or create Workspace Customer by email " +
							"perhaps a failed query or mapping",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				if err != nil {
					slog.Error(
						"failed to get or create Workspace Customer by email" +
							"something went wrong",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			case "phone":
				if !customer.Phone.Valid {
					slog.Error("phone is required for createBy phone")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}

				slog.Info("create Customer by phone")
				customer, isCreated, err = customer.GetOrCreateWrkCustomerByPhone(ctx, db)

				if errors.Is(err, model.ErrEmpty) {
					slog.Warn(
						"Customer not found or does not exist after creating" +
							"perhaps the customer is not created or is not returned",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				if errors.Is(err, model.ErrQuery) {
					slog.Error(
						"failed to get or create Workspace Customer by phone " +
							"perhaps a failed query or mapping",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				if err != nil {
					slog.Error(
						"failed to get or create Workspace Customer by phone " +
							"something went wrong",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			case "externalId":
				if !customer.ExternalId.Valid {
					slog.Error("externalId is required for createBy externalId")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}

				slog.Info("create Customer by externalId")
				customer, isCreated, err = customer.GetOrCreateWrkCustomerByExtId(ctx, db)

				if errors.Is(err, model.ErrEmpty) {
					slog.Warn(
						"Customer not found or does not exist after creating" +
							"perhaps the customer is not created or is not returned",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				if errors.Is(err, model.ErrQuery) {
					slog.Error(
						"failed to get or create Workspace Customer by externalId " +
							"perhaps a failed query or mapping",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				if err != nil {
					slog.Error(
						"failed to get or create Workspace Customer by externalId" +
							"something went wrong",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			default:
				slog.Warn("unsupported createBy value")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		} else {
			slog.Info("based on identifiers check for Customer in Workspace", slog.String("workspaceId", workspaceId))
			if customer.ExternalId.Valid {
				slog.Info("get customer by externalId")
				customer, err = customer.GetWrkCustomerByExtId(ctx, db)
				if errors.Is(err, model.ErrEmpty) {
					slog.Warn(
						"Customer not found by externalId" +
							"perhaps the customer is not created or is not returned",
					)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				if errors.Is(err, model.ErrQuery) {
					slog.Error(
						"failed to get Workspace Customer by externalId " +
							"perhaps a failed query or mapping",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				if err != nil {
					slog.Error(
						"failed to get Workspace Customer by externalId" +
							"something went wrong",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else if customer.Email.Valid {
				slog.Info("get customer by email")
				customer, err = customer.GetWrkCustomerByEmail(ctx, db)

				if errors.Is(err, model.ErrEmpty) {
					slog.Warn(
						"Customer not found by email" +
							"perhaps the customer is not created or is not returned",
					)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				if errors.Is(err, model.ErrQuery) {
					slog.Error(
						"failed to get Workspace Customer by email " +
							"perhaps a failed query or mapping",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				if err != nil {
					slog.Error(
						"failed to get Workspace Customer by email" +
							"something went wrong",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else if customer.Phone.Valid {
				slog.Info("get customer by phone")
				customer, err = customer.GetWrkCustomerByPhone(ctx, db)

				if errors.Is(err, model.ErrEmpty) {
					slog.Warn(
						"Customer not found by phone" +
							"perhaps the customer is not created or is not returned",
					)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				if errors.Is(err, model.ErrQuery) {
					slog.Error(
						"failed to get Workspace Customer by phone " +
							"perhaps a failed query or mapping",
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				if err != nil {
					slog.Error(
						"failed to get Workspace Customer by phone" +
							"something went wrong",
					)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
			} else {
				fmt.Println("unsupported customer identifier")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}

		slog.Info("got Workspace Customer or created",
			slog.String("customerId", customer.CustomerId),
			slog.Bool("isCreated", isCreated),
		)
		slog.Info("issue Customer JWT token")
		jwt, err := customer.MakeJWT()
		if err != nil {
			slog.Error("failed to make jwt token with error", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resp = CustomerTIRespPayload{
			Create:     isCreated,
			CustomerId: customer.CustomerId,
			Jwt:        jwt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(
				"failed to encode response to json "+
					"might need to check the json encoding defn",
				slog.String("customerId", customer.CustomerId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	})
}

func handleInitCustomerThreadChat(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var message ThChatReqPayload

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		customer, err := AuthenticateCustomer(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspace := model.Workspace{WorkspaceId: customer.WorkspaceId}
		workspace, err = workspace.GetById(ctx, db)
		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"workspace not found or does not exist for customer",
				slog.String("workspaceId", customer.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get workspace by id "+
					"perhaps a failed query or mapping",
				slog.String("workspaceId", customer.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err != nil {
			slog.Error(
				"failed to get workspace by id "+
					"something went wrong",
				slog.String("workspaceId", customer.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// TODO: fix to use struct with values.
		th, thm, err := model.ThreadChat{}.CreateCustomerThChat(ctx, db, workspace, customer, message.Message)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"failed to create thread chat for customer",
				slog.String("customerId", customer.CustomerId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to create thread chat for customer "+
					"perhaps a failed query or mapping",
				slog.String("customerId", customer.CustomerId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to create thread chat for customer "+
					"something went wrong",
				slog.String("customerId", customer.CustomerId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		messages := make([]ThChatMessageRespPayload, 0, 1)

		var msgCustomerRepr *ThCustomerRespPayload
		var msgMemberRepr *ThMemberRespPayload

		// for thread message - either of them
		if thm.CustomerId.Valid {
			msgCustomerRepr = &ThCustomerRespPayload{
				CustomerId: thm.CustomerId.String,
				Name:       thm.CustomerName,
			}
		} else if thm.MemberId.Valid {
			msgMemberRepr = &ThMemberRespPayload{
				MemberId: thm.MemberId.String,
				Name:     thm.MemberName,
			}
		}

		threadMessage := ThChatMessageRespPayload{
			ThreadChatId:        th.ThreadChatId,
			ThreadChatMessageId: thm.ThreadChatMessageId,
			Body:                thm.Body,
			Sequence:            thm.Sequence,
			Customer:            msgCustomerRepr,
			Member:              msgMemberRepr,
			CreatedAt:           thm.CreatedAt,
			UpdatedAt:           thm.UpdatedAt,
		}

		messages = append(messages, threadMessage)

		var threadAssigneeRepr *ThMemberRespPayload

		// for thread
		threadCustomerRepr := ThCustomerRespPayload{
			CustomerId: th.CustomerId,
			Name:       th.CustomerName,
		}

		// for thread
		if th.AssigneeId.Valid {
			threadAssigneeRepr = &ThMemberRespPayload{
				MemberId: th.AssigneeId.String,
				Name:     th.AssigneeName,
			}
		}

		resp := ThChatRespPayload{
			ThreadId:  th.ThreadChatId,
			Sequence:  th.Sequence,
			Status:    th.Status,
			Customer:  threadCustomerRepr,
			Assignee:  threadAssigneeRepr,
			CreatedAt: th.CreatedAt,
			UpdatedAt: th.UpdatedAt,
			Messages:  messages,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(
				"failed to encode thread chat to json "+
					"might need to check the json encoding defn",
				slog.String("threadChatId", th.ThreadChatId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetCustomerThreadChats(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customer, err := AuthenticateCustomer(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspace := model.Workspace{WorkspaceId: customer.WorkspaceId}
		workspace, err = workspace.GetById(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"workspace not found or does not exist for customer",
				slog.String("workspaceId", customer.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get workspace by id "+
					"perhaps a failed query or mapping",
				slog.String("workspaceId", customer.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err != nil {
			slog.Error(
				"failed to get workspace by id "+
					"something went wrong",
				slog.String("workspaceId", customer.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		th := model.ThreadChat{WorkspaceId: workspace.WorkspaceId, CustomerId: customer.CustomerId}
		ths, err := th.GetListByWorkspaceCustomerId(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"no thread chats found for customer",
				slog.String("customerId", customer.CustomerId),
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(ths); err != nil {
				slog.Error(
					"failed to encode thread chats to json "+
						"might need to check the json encoding defn",
					slog.String("customerId", customer.CustomerId),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get thread chats for customer "+
					"perhaps a failed query or mapping",
				slog.String("customerId", customer.CustomerId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get list of thread chats for customer "+
					"something went wrong",
				slog.String("customerId", customer.CustomerId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		threads := make([]ThChatRespPayload, 0, 100)
		for _, th := range ths {
			messages := make([]ThChatMessageRespPayload, 0, 1)

			var threadAssigneeRepr *ThMemberRespPayload
			var msgCustomerRepr *ThCustomerRespPayload
			var msgMemberRepr *ThMemberRespPayload

			// for thread
			threadCustomerRepr := ThCustomerRespPayload{
				CustomerId: th.ThreadChat.CustomerId,
				Name:       th.ThreadChat.CustomerName,
			}

			// for thread
			if th.ThreadChat.AssigneeId.Valid {
				threadAssigneeRepr = &ThMemberRespPayload{
					MemberId: th.ThreadChat.AssigneeId.String,
					Name:     th.ThreadChat.AssigneeName,
				}
			}

			// for thread message - either of them
			if th.Message.CustomerId.Valid {
				msgCustomerRepr = &ThCustomerRespPayload{
					CustomerId: th.Message.CustomerId.String,
					Name:       th.Message.CustomerName,
				}
			} else if th.Message.MemberId.Valid {
				msgMemberRepr = &ThMemberRespPayload{
					MemberId: th.Message.MemberId.String,
					Name:     th.Message.MemberName,
				}
			}

			message := ThChatMessageRespPayload{
				ThreadChatId:        th.ThreadChat.ThreadChatId,
				ThreadChatMessageId: th.Message.ThreadChatMessageId,
				Body:                th.Message.Body,
				Sequence:            th.Message.Sequence,
				Customer:            msgCustomerRepr,
				Member:              msgMemberRepr,
				CreatedAt:           th.Message.CreatedAt,
				UpdatedAt:           th.Message.UpdatedAt,
			}
			messages = append(messages, message)
			threads = append(threads, ThChatRespPayload{
				ThreadId:  th.ThreadChat.ThreadChatId,
				Sequence:  th.ThreadChat.Sequence,
				Status:    th.ThreadChat.Status,
				Customer:  threadCustomerRepr,
				Assignee:  threadAssigneeRepr,
				CreatedAt: th.ThreadChat.CreatedAt,
				UpdatedAt: th.ThreadChat.UpdatedAt,
				Messages:  messages,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(threads); err != nil {
			slog.Error(
				"failed to encode thread chats to json "+
					"might need to check the json encoding defn",
				slog.String("customerId", customer.CustomerId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleCreateCustomerThChatMessage(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		threadId := r.PathValue("threadId")

		var message ThChatReqPayload

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		customer, err := AuthenticateCustomer(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspace := model.Workspace{WorkspaceId: customer.WorkspaceId}
		_, err = workspace.GetById(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"workspace not found or does not exist for customer",
				slog.String("workspaceId", customer.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get workspace by id "+
					"perhaps a failed query or mapping",
				slog.String("workspaceId", customer.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get workspace by id "+
					"something went wrong",
				slog.String("workspaceId", customer.WorkspaceId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		th := model.ThreadChat{ThreadChatId: threadId}
		th, err = th.GetById(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"thread chat not found or does not exist for customer",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get thread chat by id "+
					"perhaps a failed query or mapping",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get thread chat by id "+
					"something went wrong",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		thm := model.ThreadChatMessage{ThreadChatId: th.ThreadChatId}
		thm, err = thm.CreateCustomerThChatMessage(ctx, db, customer, message.Message)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"no thread chat message found for customer after creation",
				slog.String("threadChatId", th.ThreadChatId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to create thread chat message for customer "+
					"perhaps a failed query or mapping",
				slog.String("threadChatId", th.ThreadChatId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to create thread chat message for customer "+
					"something went wrong",
				slog.String("threadChatId", th.ThreadChatId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var threadAssigneeRepr *ThMemberRespPayload
		var msgCustomerRepr *ThCustomerRespPayload
		var msgMemberRepr *ThMemberRespPayload

		// for thread
		threadCustomerRepr := ThCustomerRespPayload{
			CustomerId: th.CustomerId,
			Name:       th.CustomerName,
		}

		// for thread
		if th.AssigneeId.Valid {
			threadAssigneeRepr = &ThMemberRespPayload{
				MemberId: th.AssigneeId.String,
				Name:     th.AssigneeName,
			}
		}

		// for thread message - either of them
		if thm.CustomerId.Valid {
			msgCustomerRepr = &ThCustomerRespPayload{
				CustomerId: thm.CustomerId.String,
				Name:       thm.CustomerName,
			}
		} else if thm.MemberId.Valid {
			msgMemberRepr = &ThMemberRespPayload{
				MemberId: thm.MemberId.String,
				Name:     thm.MemberName,
			}
		}

		threadMessage := ThChatMessageRespPayload{
			ThreadChatId:        th.ThreadChatId,
			ThreadChatMessageId: thm.ThreadChatMessageId,
			Body:                thm.Body,
			Sequence:            thm.Sequence,
			Customer:            msgCustomerRepr,
			Member:              msgMemberRepr,
			CreatedAt:           thm.CreatedAt,
			UpdatedAt:           thm.UpdatedAt,
		}

		messages := make([]ThChatMessageRespPayload, 0, 1)
		messages = append(messages, threadMessage)
		resp := ThChatRespPayload{
			ThreadId:  th.ThreadChatId,
			Sequence:  th.Sequence,
			Status:    th.Status,
			Customer:  threadCustomerRepr,
			Assignee:  threadAssigneeRepr,
			CreatedAt: th.CreatedAt,
			UpdatedAt: th.UpdatedAt,
			Messages:  messages,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(
				"failed to encode thread chat message to json "+
					"might need to check the json encoding defn",
				slog.String("threadChatId", th.ThreadChatId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetCustomerThChatMessages(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := AuthenticateCustomer(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		threadId := r.PathValue("threadId")
		th := model.ThreadChat{ThreadChatId: threadId}

		th, err = th.GetById(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"thread chat not found or does not exist for customer",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get thread chat by id "+
					"perhaps a failed query or mapping",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get thread chat by id "+
					"something went wrong",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		thc := model.ThreadChatMessage{ThreadChatId: th.ThreadChatId}
		results, err := thc.GetListByThreadChatId(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"no thread chat messages found for customer",
				slog.String("threadChatId", th.ThreadChatId),
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(results); err != nil {
				slog.Error(
					"failed to encode thread chat messages to json "+
						"might need to check the json encoding defn",
					slog.String("threadChatId", th.ThreadChatId),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get thread chat messages for customer "+
					"perhaps a failed query or mapping",
				slog.String("threadChatId", th.ThreadChatId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get list of thread chat messages for customer "+
					"something went wrong",
				slog.String("threadChatId", th.ThreadChatId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		messages := make([]ThChatMessageRespPayload, 0, 100)
		for _, thm := range results {
			var msgCustomerRepr *ThCustomerRespPayload
			var msgMemberRepr *ThMemberRespPayload

			// for thread message - either of them
			if thm.CustomerId.Valid {
				msgCustomerRepr = &ThCustomerRespPayload{
					CustomerId: thm.CustomerId.String,
					Name:       thm.CustomerName,
				}
			} else if thm.MemberId.Valid {
				msgMemberRepr = &ThMemberRespPayload{
					MemberId: thm.MemberId.String,
					Name:     thm.MemberName,
				}
			}

			threadMessage := ThChatMessageRespPayload{
				ThreadChatId:        th.ThreadChatId,
				ThreadChatMessageId: thm.ThreadChatMessageId,
				Body:                thm.Body,
				Sequence:            thm.Sequence,
				Customer:            msgCustomerRepr,
				Member:              msgMemberRepr,
				CreatedAt:           thm.CreatedAt,
				UpdatedAt:           thm.UpdatedAt,
			}

			messages = append(messages, threadMessage)
		}

		var threadAssigneeRepr *ThMemberRespPayload

		// for thread
		threadCustomerRepr := ThCustomerRespPayload{
			CustomerId: th.CustomerId,
			Name:       th.CustomerName,
		}

		// for thread
		if th.AssigneeId.Valid {
			threadAssigneeRepr = &ThMemberRespPayload{
				MemberId: th.AssigneeId.String,
				Name:     th.AssigneeName,
			}
		}

		resp := ThChatRespPayload{
			ThreadId:  th.ThreadChatId,
			Sequence:  th.Sequence,
			Status:    th.Status,
			Customer:  threadCustomerRepr,
			Assignee:  threadAssigneeRepr,
			CreatedAt: th.CreatedAt,
			UpdatedAt: th.UpdatedAt,
			Messages:  messages,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(
				"failed to encode thread chat messages to json "+
					"might need to check the json encoding defn",
				slog.String("threadChatId", th.ThreadChatId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleCreateMemberThChatMessage(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		workspaceId := r.PathValue("workspaceId")
		threadId := r.PathValue("threadId")

		var message ThChatReqPayload

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		account, err := AuthenticateAccount(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		member := model.Member{WorkspaceId: workspaceId, AccountId: account.AccountId}
		member, err = member.GetWorkspaceMemberByAccountId(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"no member found for account",
				slog.String("accountId", account.AccountId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get member "+
					"perhaps a failed query or mapping",
				slog.String("accountId", account.AccountId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get member "+
					"something went wrong",
				slog.String("accountId", account.AccountId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		th := model.ThreadChat{ThreadChatId: threadId}

		th, err = th.GetById(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"no thread chat found",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get thread chat "+
					"perhaps a failed query or mapping",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get thread chat "+
					"something went wrong",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		thm := model.ThreadChatMessage{ThreadChatId: th.ThreadChatId}
		thm, err = thm.CreateMemberThChatMessage(ctx, db, member, message.Message)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"no thread chat message found after creation",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to create thread chat message "+
					"perhaps a failed query or mapping",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err != nil {
			slog.Error(
				"failed to create thread chat message "+
					"something went wrong",
				slog.String("threadChatId", threadId),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !th.AssigneeId.Valid {
			slog.Info("Thread Chat not yet assigned will assign Member...")
			thAssigned := th // make a temp copy before assigning
			thAssigned.AssigneeId = NullString(&member.MemberId)
			thAssigned, err := thAssigned.AssignMember(ctx, db)
			if err != nil {
				slog.Error("(silent) failed to assign Member to Thread Chat", slog.Any("error", err))
			} else {
				th = thAssigned // update the original with assigned
			}
		}

		var threadAssigneeRepr *ThMemberRespPayload

		var msgCustomerRepr *ThCustomerRespPayload
		var msgMemberRepr *ThMemberRespPayload

		// for thread
		threadCustomerRepr := ThCustomerRespPayload{
			CustomerId: th.CustomerId,
			Name:       th.CustomerName,
		}

		// for thread
		if th.AssigneeId.Valid {
			threadAssigneeRepr = &ThMemberRespPayload{
				MemberId: th.AssigneeId.String,
				Name:     th.AssigneeName,
			}
		}

		// for thread message - either of them
		if thm.CustomerId.Valid {
			msgCustomerRepr = &ThCustomerRespPayload{
				CustomerId: thm.CustomerId.String,
				Name:       thm.CustomerName,
			}
		} else if thm.MemberId.Valid {
			msgMemberRepr = &ThMemberRespPayload{
				MemberId: thm.MemberId.String,
				Name:     thm.MemberName,
			}
		}

		threadMessage := ThChatMessageRespPayload{
			ThreadChatId:        th.ThreadChatId,
			ThreadChatMessageId: thm.ThreadChatMessageId,
			Body:                thm.Body,
			Sequence:            thm.Sequence,
			Customer:            msgCustomerRepr,
			Member:              msgMemberRepr,
			CreatedAt:           thm.CreatedAt,
			UpdatedAt:           thm.UpdatedAt,
		}

		messages := make([]ThChatMessageRespPayload, 0, 1)
		messages = append(messages, threadMessage)
		resp := ThChatRespPayload{
			ThreadId:  th.ThreadChatId,
			Sequence:  th.Sequence,
			Status:    th.Status,
			Customer:  threadCustomerRepr,
			Assignee:  threadAssigneeRepr,
			CreatedAt: th.CreatedAt,
			UpdatedAt: th.UpdatedAt,
			Messages:  messages,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error(
				"failed to encode response",
				slog.Any("error", err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

// TODO: later
func handleInitCustomerThreadQA(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var query ThreadQAReqPayload

		err := json.NewDecoder(r.Body).Decode(&query)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		customer, err := AuthenticateCustomer(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspace, err := model.Workspace{WorkspaceId: customer.WorkspaceId}.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		tq := model.ThreadQA{
			WorkspaceId: workspace.WorkspaceId,
			CustomerId:  customer.CustomerId,
			Query:       query.Query,
		}

		tq, err = tq.Create(ctx, db)
		if err != nil {
			slog.Error("failed to create query", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		reqId := xid.New()
		wrkLLM := LLM{WorkspaceId: workspace.WorkspaceId, Prompt: tq.Query, RequestId: reqId.String()}
		llmr, err := wrkLLM.Generate()
		if err != nil {
			slog.Error("failed to generate llm response", "error", err)
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		answerId := xid.New()
		tqa := model.ThreadQAA{
			WorkspaceId: workspace.WorkspaceId,
			ThreadQAId:  tq.ThreadId,
			AnswerId:    answerId.String(),
			Answer:      llmr.Text,
		}

		tqa, err = tqa.Create(ctx, db)
		if err != nil {
			slog.Error("failed to create thread question answer", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ans := make([]ThreadQAARespPayload, 0, 1)
		ans = append(ans, ThreadQAARespPayload{
			AnswerId: tqa.AnswerId,
			Answer:   tqa.Answer,
			Eval:     tqa.Eval,
			Sequence: tqa.Sequence,
		})
		resp := ThreadQARespPayload{
			ThreadId:  tq.ThreadId,
			Query:     tq.Query,
			Sequence:  tq.Sequence,
			CreatedAt: tq.CreatedAt,
			UpdatedAt: tq.UpdatedAt,
			Answers:   ans,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetThreadChats(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, err := AuthenticateAccount(ctx, db, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspaceId := r.PathValue("workspaceId")
		workspace := model.Workspace{AccountId: account.AccountId, WorkspaceId: workspaceId}
		workspace, err = workspace.GetAccountWorkspace(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"workspace not found or does not exist for account",
				"accountId", account.AccountId, "workspaceId", workspaceId,
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get workspace by id "+
					"perhaps a failed query or mapping",
				"accountId", account.AccountId, "workspaceId", workspaceId,
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get workspace by id "+
					"something went wrong",
				"accountId", account.AccountId, "workspaceId", workspaceId,
			)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		th := model.ThreadChat{WorkspaceId: workspace.WorkspaceId}
		ths, err := th.GetListByWorkspace(ctx, db)

		if errors.Is(err, model.ErrEmpty) {
			slog.Warn(
				"no thread chats found for workspace",
				"workspaceId", workspace.WorkspaceId,
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(ths); err != nil {
				slog.Error(
					"failed to encode thread chats to json "+
						"might need to check the json encoding defn",
					"workspaceId", workspace.WorkspaceId,
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		if errors.Is(err, model.ErrQuery) {
			slog.Error(
				"failed to get thread chats for workspace "+
					"perhaps a failed query or mapping",
				"workspaceId", workspace.WorkspaceId,
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err != nil {
			slog.Error(
				"failed to get list of thread chats for workspace "+
					"something went wrong",
				"workspaceId", workspace.WorkspaceId,
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		threads := make([]ThChatRespPayload, 0, 100)
		for _, th := range ths {
			messages := make([]ThChatMessageRespPayload, 0, 1)

			var threadAssigneeRepr *ThMemberRespPayload
			var msgCustomerRepr *ThCustomerRespPayload
			var msgMemberRepr *ThMemberRespPayload

			// for thread
			threadCustomerRepr := ThCustomerRespPayload{
				CustomerId: th.ThreadChat.CustomerId,
				Name:       th.ThreadChat.CustomerName,
			}

			// for thread
			if th.ThreadChat.AssigneeId.Valid {
				threadAssigneeRepr = &ThMemberRespPayload{
					MemberId: th.ThreadChat.AssigneeId.String,
					Name:     th.ThreadChat.AssigneeName,
				}
			}

			// for thread message - either of them
			if th.Message.CustomerId.Valid {
				msgCustomerRepr = &ThCustomerRespPayload{
					CustomerId: th.Message.CustomerId.String,
					Name:       th.Message.CustomerName,
				}
			} else if th.Message.MemberId.Valid {
				msgMemberRepr = &ThMemberRespPayload{
					MemberId: th.Message.MemberId.String,
					Name:     th.Message.MemberName,
				}
			}

			message := ThChatMessageRespPayload{
				ThreadChatId:        th.ThreadChat.ThreadChatId,
				ThreadChatMessageId: th.Message.ThreadChatMessageId,
				Body:                th.Message.Body,
				Sequence:            th.Message.Sequence,
				Customer:            msgCustomerRepr,
				Member:              msgMemberRepr,
				CreatedAt:           th.Message.CreatedAt,
				UpdatedAt:           th.Message.UpdatedAt,
			}
			messages = append(messages, message)
			threads = append(threads, ThChatRespPayload{
				ThreadId:  th.ThreadChat.ThreadChatId,
				Sequence:  th.ThreadChat.Sequence,
				Status:    th.ThreadChat.Status,
				Customer:  threadCustomerRepr,
				Assignee:  threadAssigneeRepr,
				CreatedAt: th.ThreadChat.CreatedAt,
				UpdatedAt: th.ThreadChat.UpdatedAt,
				Messages:  messages,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(threads); err != nil {
			slog.Error(
				"failed to encode thread chats to json "+
					"might need to check the json encoding defn",
				"workspaceId", workspace.WorkspaceId,
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})

}
