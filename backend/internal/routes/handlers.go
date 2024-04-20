package routes

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/internal/auth"
	"github.com/zyghq/zyg/internal/model"
)

const DefaultAuthProvider string = "supabase"

func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// request
type PATReqPayload struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// request
type WorkspaceReqPayload struct {
	Name string `json:"name"`
}

// request
type CustomerTIReqPayload struct {
	Create   bool    `json:"create"`
	CreateBy *string `json:"createBy"` // optional
	Customer struct {
		ExternalId *string `json:"externalId"` // optional
		Email      *string `json:"email"`      // optional
		Phone      *string `json:"phone"`      // optional
	} `json:"customer"`
}

// response
type CustomerTIRespPayload struct {
	Create     bool   `json:"create"`
	CustomerId string `json:"customerId"`
	Jwt        string `json:"jwt"`
}

// request
type ThChatReqPayload struct {
	Message string `json:"message"`
}

// response
type ThCustomerRespPayload struct {
	CustomerId string
	Name       sql.NullString
}

// response
type ThMemberRespPayload struct {
	MemberId string
	Name     sql.NullString
}

func (m ThMemberRespPayload) MarshalJSON() ([]byte, error) {
	var name *string
	// if m.MemberId.Valid {
	// 	memberId = &m.MemberId.String
	// }
	if m.Name.Valid {
		name = &m.Name.String
	}
	aux := &struct {
		MemberId string  `json:"memberId"`
		Name     *string `json:"name"`
	}{
		MemberId: m.MemberId,
		Name:     name,
	}
	return json.Marshal(aux)
}

func (c ThCustomerRespPayload) MarshalJSON() ([]byte, error) {
	var name *string
	if c.Name.Valid {
		name = &c.Name.String
	}
	aux := &struct {
		CustomerId string  `json:"customerId"`
		Name       *string `json:"name"`
	}{
		CustomerId: c.CustomerId,
		Name:       name,
	}
	return json.Marshal(aux)
}

// response
type ThChatMessageRespPayload struct {
	ThreadChatId        string
	ThreadChatMessageId string
	Body                string
	Sequence            int
	Customer            *ThCustomerRespPayload
	Member              *ThMemberRespPayload
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (thcresp ThChatMessageRespPayload) MarshalJSON() ([]byte, error) {
	var customer *ThCustomerRespPayload
	var member *ThMemberRespPayload

	if thcresp.Customer != nil {
		customer = thcresp.Customer
	}

	if thcresp.Member != nil {
		member = thcresp.Member
	}

	aux := &struct {
		ThreadChatId        string                 `json:"threadChatId"`
		ThreadChatMessageId string                 `json:"threadChatMessageId"`
		Body                string                 `json:"body"`
		Sequence            int                    `json:"sequence"`
		Customer            *ThCustomerRespPayload `json:"customer,omitempty"`
		Member              *ThMemberRespPayload   `json:"member,omitempty"`
		CreatedAt           string                 `json:"createdAt"`
		UpdatedAt           string                 `json:"updatedAt"`
	}{
		ThreadChatId:        thcresp.ThreadChatId,
		ThreadChatMessageId: thcresp.ThreadChatMessageId,
		Body:                thcresp.Body,
		Sequence:            thcresp.Sequence,
		Customer:            customer,
		Member:              member,
		CreatedAt:           thcresp.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           thcresp.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

// response
type ThChatRespPayload struct {
	ThreadId  string
	Sequence  int
	Status    string
	Customer  ThCustomerRespPayload
	Assignee  *ThMemberRespPayload
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []ThChatMessageRespPayload
}

func (thresp ThChatRespPayload) MarshalJSON() ([]byte, error) {
	var assignee *ThMemberRespPayload

	if thresp.Assignee != nil {
		assignee = thresp.Assignee
	}

	aux := &struct {
		ThreadId  string                     `json:"threadId"`
		Sequence  int                        `json:"sequence"`
		Status    string                     `json:"status"`
		Customer  ThCustomerRespPayload      `json:"customer"`
		Assignee  *ThMemberRespPayload       `json:"assignee"`
		CreatedAt string                     `json:"createdAt"`
		UpdatedAt string                     `json:"updatedAt"`
		Messages  []ThChatMessageRespPayload `json:"messages"`
	}{
		ThreadId:  thresp.ThreadId,
		Sequence:  thresp.Sequence,
		Status:    thresp.Status,
		Customer:  thresp.Customer,
		Assignee:  assignee,
		CreatedAt: thresp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: thresp.UpdatedAt.Format(time.RFC3339),
		Messages:  thresp.Messages,
	}
	return json.Marshal(aux)
}

// request
type ThreadQAReqPayload struct {
	Query string `json:"query"`
}

// response
type ThreadQAARespPayload struct {
	AnswerId string
	Answer   string
	Eval     sql.NullInt32
	Sequence int
}

func (tqar ThreadQAARespPayload) MarshalJSON() ([]byte, error) {
	var eval *int32
	if tqar.Eval.Valid {
		eval = &tqar.Eval.Int32
	}
	aux := &struct {
		AnswerId string `json:"answerId"`
		Answer   string `json:"answer"`
		Eval     *int32 `json:"eval"`
		Sequence int    `json:"sequence"`
	}{
		AnswerId: tqar.AnswerId,
		Answer:   tqar.Answer,
		Eval:     eval,
		Sequence: tqar.Sequence,
	}
	return json.Marshal(aux)
}

type ThreadQARespPayload struct {
	ThreadId  string
	Query     string
	Sequence  int
	CreatedAt time.Time
	UpdatedAt time.Time
	Answers   []ThreadQAARespPayload
}

func (thr ThreadQARespPayload) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ThreadId  string                 `json:"threadId"`
		Query     string                 `json:"query"`
		Sequence  int                    `json:"sequence"`
		CreatedAt string                 `json:"createdAt"`
		UpdatedAt string                 `json:"updatedAt"`
		Answers   []ThreadQAARespPayload `json:"answers"`
	}{
		ThreadId:  thr.ThreadId,
		Query:     thr.Query,
		Sequence:  thr.Sequence,
		CreatedAt: thr.CreatedAt.Format(time.RFC3339),
		UpdatedAt: thr.UpdatedAt.Format(time.RFC3339),
		Answers:   thr.Answers,
	}
	return json.Marshal(aux)
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

func AuthenticateAccount(ctx context.Context, db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) (model.Account, error) {
	var account model.Account

	ath := r.Header.Get("Authorization")
	if ath == "" {
		return account, fmt.Errorf("cannot authenticate without authorization header")
	}

	cred := strings.Split(ath, " ")
	scheme := strings.ToLower(cred[0])

	if scheme == "token" {
		fmt.Println("authenticate with PAT")
		token := model.AccountPAT{Token: cred[1]}
		account, err := token.GetByToken(ctx, db)
		if err != nil {
			return account, fmt.Errorf("failed to authenticate with error: %v", err)
		}
		fmt.Printf("authenticated account with accountId: %s\n", account.AccountId)
		return account, nil
	} else if scheme == "bearer" {
		fmt.Println("authenticate with JWTs")
		hmacSecret, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			return account, fmt.Errorf("failed to get env SUPABASE_JWT_SECRET with error: %v", err)
		}
		ac, err := auth.ParseJWTToken(cred[1], []byte(hmacSecret))
		if err != nil {
			return account, fmt.Errorf("failed to parse JWT token with error: %v", err)
		}
		sub, err := ac.RegisteredClaims.GetSubject()
		if err != nil {
			return account, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}
		fmt.Printf("authenticated account with auth user id: %s\n", sub)

		// fetch the authenticated account
		account = model.Account{AuthUserId: sub}
		account, err = account.GetByAuthUserId(ctx, db)
		if err != nil {
			return account, fmt.Errorf("failed to get account by auth user id: %s with error: %v", sub, err)
		}
		// return the authenticated account
		return account, nil
	} else {
		return account, fmt.Errorf("unsupported scheme: `%s` cannot authenticate", scheme)
	}
}

func AuthenticateCustomer(ctx context.Context, db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) (model.Customer, error) {
	var customer model.Customer

	ath := r.Header.Get("Authorization")
	if ath == "" {
		return customer, fmt.Errorf("cannot authenticate without authorization header")
	}

	cred := strings.Split(ath, " ")
	scheme := strings.ToLower(cred[0])

	if scheme == "bearer" {
		fmt.Println("authenticate with JWTs")
		hmacSecret, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			return customer, fmt.Errorf("failed to get env SUPABASE_JWT_SECRET with error: %v", err)
		}
		cc, err := auth.ParseCustomerJWTToken(cred[1], []byte(hmacSecret))
		if err != nil {
			return customer, fmt.Errorf("failed to parse JWT token with error: %v", err)
		}
		sub, err := cc.RegisteredClaims.GetSubject()
		if err != nil {
			return customer, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}
		fmt.Printf("authenticated customer with id: %s\n", sub)

		// fetch the authenticated customer
		customer = model.Customer{WorkspaceId: cc.WorkspaceId, CustomerId: sub}
		customer, err = customer.GetWrkCustomerById(ctx, db)
		if err != nil {
			return customer, fmt.Errorf("failed to get customer by customer id: %s with error: %v", customer.CustomerId, err)
		}
		// return the authenticated customer
		return customer, nil
	} else {
		return customer, fmt.Errorf("unsupported scheme: `%s` cannot authenticate", scheme)
	}
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

		ath := r.Header.Get("Authorization")
		if ath == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		cred := strings.Split(ath, " ")
		scheme := strings.ToLower(cred[0])

		if scheme == "token" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if scheme == "bearer" {
			hmacSecret, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			ac, err := auth.ParseJWTToken(cred[1], []byte(hmacSecret))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			sub, err := ac.RegisteredClaims.GetSubject()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			account := model.Account{AuthUserId: sub, Email: ac.Email, Provider: DefaultAuthProvider}
			account, isCreated, err := account.GetOrCreateByAuthUserId(ctx, db)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if isCreated {
				fmt.Printf("account created with accountId: %s\n", account.AccountId)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				if err := json.NewEncoder(w).Encode(account); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(account); err != nil {
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

		account, err := AuthenticateAccount(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ap := model.AccountPAT{
			AccountId: account.AccountId,
			Name:      rb.Name,
			UnMask:    true, // unmask only once created
		}
		ap.Description = NullString(rb.Description)

		ap, err = ap.Create(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(ap); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetPATs(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, err := AuthenticateAccount(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ap := model.AccountPAT{AccountId: account.AccountId}
		aps, err := ap.GetListByAccountId(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(aps); err != nil {
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

		account, err := AuthenticateAccount(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		workspace := model.Workspace{AccountId: account.AccountId, Name: rb.Name}
		workspace, err = workspace.Create(ctx, db)
		if err != nil {
			fmt.Printf("failed to create workspace with error: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(workspace); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetWorkspaces(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		account, err := AuthenticateAccount(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(ctx, `SELECT
			workspace_id, account_id,
			name, created_at, updated_at
			FROM workspace WHERE account_id = $1
			ORDER BY created_at
			DESC LIMIT 100`, account.AccountId)
		if err != nil {
			log.Printf("error: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		workspaces := make([]model.Workspace, 0)
		for rows.Next() {
			var workspace model.Workspace
			err = rows.Scan(
				&workspace.WorkspaceId, &workspace.AccountId,
				&workspace.Name,
				&workspace.CreatedAt, &workspace.UpdatedAt,
			)
			if err != nil {
				log.Printf("error: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			workspaces = append(workspaces, workspace)
		}

		if err := rows.Err(); err != nil {
			log.Printf("error: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(workspaces); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetCustomer(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customer, err := AuthenticateCustomer(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(customer); err != nil {
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
			fmt.Printf("failed to decode request body error: %v\n", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		account, err := AuthenticateAccount(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspaceId := r.PathValue("workspaceId")
		fmt.Printf("issue token for customer in workspaceId: %v\n", workspaceId)

		tw := model.Workspace{WorkspaceId: workspaceId, AccountId: account.AccountId}
		workspace, err := tw.GetAccountWorkspace(ctx, db)
		if err != nil {
			fmt.Printf("failed to get account workspace or does not exist with error: %v\n", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		customer := model.Customer{
			WorkspaceId: workspace.WorkspaceId,
		}
		customer.ExternalId = NullString(rb.Customer.ExternalId)
		customer.Email = NullString(rb.Customer.Email)
		customer.Phone = NullString(rb.Customer.Phone)
		if !customer.ExternalId.Valid && !customer.Email.Valid && !customer.Phone.Valid {
			fmt.Println("at least one of `externalId`, `email` or `phone` is required")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if rb.Create {
			if rb.CreateBy == nil {
				fmt.Println("requires `createBy` when `create` is enabled")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			createBy := *rb.CreateBy
			fmt.Printf("create the customer if does not exists by %s\n", createBy)
			switch createBy {
			case "email":
				if !customer.Email.Valid {
					fmt.Println("`email` is required for `createBy` email")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				email := customer.Email.String
				fmt.Printf("create the customer by email %s\n", email)
				customer, isCreated, err := customer.GetOrCreateWrkCustomerByEmail(ctx, db)
				if err != nil {
					fmt.Printf("failed to get or create customer by email %s with error: %v\n", email, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				fmt.Printf("customer id: %s is created: %v\n", customer.CustomerId, isCreated)
				jwt, err := customer.MakeJWT()
				if err != nil {
					fmt.Printf("failed to make jwt token with error: %v\n", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				resp := CustomerTIRespPayload{
					Create:     isCreated,
					CustomerId: customer.CustomerId,
					Jwt:        jwt,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				return
			case "phone":
				if !customer.Phone.Valid {
					fmt.Println("`phone` is required for `createBy` phone")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				phone := customer.Phone.String
				fmt.Printf("create the customer by phone %s\n", phone)
				customer, isCreated, err := customer.GetOrCreateWrkCustomerByPhone(ctx, db)
				if err != nil {
					fmt.Printf("failed to get or create customer by phone %s with error: %v\n", phone, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				fmt.Printf("customerId: %s is created: %v\n", customer.CustomerId, isCreated)
				jwt, err := customer.MakeJWT()
				if err != nil {
					fmt.Printf("failed to make jwt token with error: %v\n", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				resp := CustomerTIRespPayload{
					Create:     isCreated,
					CustomerId: customer.CustomerId,
					Jwt:        jwt,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			case "externalId":
				if !customer.ExternalId.Valid {
					fmt.Println("`externalId` is required for `createBy` externalId")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				extId := customer.ExternalId.String
				fmt.Printf("create the customer by externalId %s\n", extId)
				customer, isCreated, err := customer.GetOrCreateWrkCustomerByExtId(ctx, db)
				if err != nil {
					fmt.Printf("failed to get or create customer by externalId %s with error: %v\n", extId, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				fmt.Printf("customerId: %s is created: %v\n", customer.CustomerId, isCreated)
				jwt, err := customer.MakeJWT()
				if err != nil {
					fmt.Printf("failed to make jwt token with error: %v\n", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				resp := CustomerTIRespPayload{
					Create:     isCreated,
					CustomerId: customer.CustomerId,
					Jwt:        jwt,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			default:
				fmt.Println("unsupported `createBy` field value")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		} else {
			fmt.Printf("based on identifiers check for customer in workspaceId: %v\n", workspaceId)
			if customer.ExternalId.Valid {
				fmt.Printf("get workspace customer by externalId %s\n", customer.ExternalId.String)
				customer, err = customer.GetWrkCustomerByExtId(ctx, db)
				if err != nil {
					// TOOD: improve logging and error handling
					fmt.Printf("failed to get customer by externalId %s with error: %v\n", customer.ExternalId.String, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				fmt.Printf("found customer with customer id: %s\n", customer.CustomerId)
				jwt, err := customer.MakeJWT()
				if err != nil {
					fmt.Printf("failed to make jwt token with error: %v\n", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				resp := CustomerTIRespPayload{
					Create:     false,
					CustomerId: customer.CustomerId,
					Jwt:        jwt,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else if customer.Email.Valid {
				fmt.Printf("get customer by email %s\n", customer.Email.String)
				customer, err = customer.GetWrkCustomerByEmail(ctx, db)
				if err != nil {
					fmt.Printf("failed to get customer by email %s with error: %v\n", customer.Email.String, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				fmt.Printf("found customer with customer id: %s\n", customer.CustomerId)
				jwt, err := customer.MakeJWT()
				if err != nil {
					fmt.Printf("failed to make jwt token with error: %v\n", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				resp := CustomerTIRespPayload{
					Create:     false,
					CustomerId: customer.CustomerId,
					Jwt:        jwt,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else if customer.Phone.Valid {
				fmt.Printf("get customer by phone %s\n", customer.Phone.String)
				customer, err = customer.GetWrkCustomerByPhone(ctx, db)
				if err != nil {
					fmt.Printf("failed to get customer by phone %s with error: %v\n", customer.Phone.String, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				fmt.Printf("found customer with customer id: %s\n", customer.CustomerId)
				jwt, err := customer.MakeJWT()
				if err != nil {
					fmt.Printf("failed to make jwt token with error: %v\n", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				resp := CustomerTIRespPayload{
					Create:     false,
					CustomerId: customer.CustomerId,
					Jwt:        jwt,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				fmt.Println("unsupported customer identifier")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
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

		customer, err := AuthenticateCustomer(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspace, err := model.Workspace{WorkspaceId: customer.WorkspaceId}.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		th, thm, err := model.ThreadChat{}.CreateCustomerThChat(ctx, db, workspace, customer, message.Message)
		if err != nil {
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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetCustomerThreadChats(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customer, err := AuthenticateCustomer(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		workspace, err := model.Workspace{WorkspaceId: customer.WorkspaceId}.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		th := model.ThreadChat{WorkspaceId: workspace.WorkspaceId, CustomerId: customer.CustomerId}
		ths, err := th.GetListByWorkspaceCustomerId(ctx, db)
		if err != nil {
			fmt.Printf("error in Get List By WorksapceId %s", err)
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

		customer, err := AuthenticateCustomer(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		_, err = model.Workspace{WorkspaceId: customer.WorkspaceId}.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		th := model.ThreadChat{ThreadChatId: threadId}
		th, err = th.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		thm := model.ThreadChatMessage{ThreadChatId: th.ThreadChatId}
		thm, err = thm.CreateCustomerThChatMessage(ctx, db, customer, message.Message)
		if err != nil {
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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleGetCustomerThChatMessages(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := AuthenticateCustomer(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		threadId := r.PathValue("threadId")
		th := model.ThreadChat{ThreadChatId: threadId}

		th, err = th.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		thc := model.ThreadChatMessage{ThreadChatId: th.ThreadChatId}
		results, err := thc.GetListByThreadChatId(ctx, db)
		if err != nil {
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

		account, err := AuthenticateAccount(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		member := model.Member{WorkspaceId: workspaceId, AccountId: account.AccountId}
		member, err = member.GetWorkspaceMemberByAccountId(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		th := model.ThreadChat{ThreadChatId: threadId}

		th, err = th.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		thm := model.ThreadChatMessage{ThreadChatId: th.ThreadChatId}
		thm, err = thm.CreateMemberThChatMessage(ctx, db, member, message.Message)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !th.AssigneeId.Valid {
			fmt.Println("Thread Chat not yet assigned will try to assign Member...")
			thAssigned := th // make a temp copy
			thAssigned.AssigneeId = NullString(&member.MemberId)
			thAssigned, err := thAssigned.AssignMember(ctx, db)
			if err != nil {
				fmt.Printf("(silent) failed to assign member to Thread Chat with error: %v\n", err)
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

		customer, err := AuthenticateCustomer(ctx, db, w, r)
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
			fmt.Printf("failed to create thread query with error: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		reqId := xid.New()
		wrkLLM := LLM{WorkspaceId: workspace.WorkspaceId, Prompt: tq.Query, RequestId: reqId.String()}
		llmr, err := wrkLLM.Generate()
		if err != nil {
			fmt.Printf("failed to generate llm response with error: %v\n", err)
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
			fmt.Printf("failed to create thread question answer with error: %v\n", err)
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
