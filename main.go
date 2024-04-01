package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"github.com/rs/xid"
)

var addr = flag.String("addr", "127.0.0.1:8080", "listen address")

func GetEnv(key string) (string, error) {
	value, status := os.LookupEnv(key)
	if !status {
		return "", fmt.Errorf("env `%s` is not set", key)
	}
	return value, nil
}

// models
type Account struct {
	AccountId  string
	Email      string
	Provider   string
	AuthUserId string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (a Account) MarshalJSON() ([]byte, error) {
	aux := &struct {
		AccountId string `json:"accountId"`
		Email     string `json:"email"`
		Provider  string `json:"provider"`
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		AccountId: a.AccountId,
		Email:     a.Email,
		Provider:  a.Provider,
		Name:      a.Name,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

// models
type AccountPat struct {
	AccountId   string
	PatId       string
	Token       string
	Name        string
	Description sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// models
type Workspace struct {
	WorkspaceId string
	AccountId   string
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (w Workspace) MarshalJSON() ([]byte, error) {
	aux := &struct {
		WorkspaceId string `json:"workspaceId"`
		AccountId   string `json:"accountId"`
		Name        string `json:"name"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		WorkspaceId: w.WorkspaceId,
		AccountId:   w.AccountId,
		Name:        w.Name,
		CreatedAt:   w.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   w.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

// models
type Customer struct {
	WorkspaceId string
	CustomerId  string
	ExternalId  sql.NullString
	Email       sql.NullString
	Phone       sql.NullString
	Name        sql.NullString
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

func (c Customer) MarshalJSON() ([]byte, error) {
	var externalId, email, phone, name *string
	if c.ExternalId.Valid {
		externalId = &c.ExternalId.String
	}
	if c.Email.Valid {
		email = &c.Email.String
	}
	if c.Phone.Valid {
		phone = &c.Phone.String
	}
	if c.Name.Valid {
		name = &c.Name.String
	}

	aux := &struct {
		WorkspaceId string  `json:"workspaceId"`
		CustomerId  string  `json:"customerId"`
		ExternalId  *string `json:"externalId"`
		Email       *string `json:"email"`
		Phone       *string `json:"phone"`
		Name        *string `json:"name"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}{
		WorkspaceId: c.WorkspaceId,
		CustomerId:  c.CustomerId,
		ExternalId:  externalId,
		Email:       email,
		Phone:       phone,
		Name:        name,
		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

// models
type LLMRRLog struct {
	WorkspaceId string
	RequestId   string
	Prompt      string
	Response    string
	Model       string
	Eval        sql.NullInt64
}

// request
type LLMRequestQuery struct {
	Q string `json:"q"`
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

// request
type CustomerTIRequest struct {
	Create   bool    `json:"create"`
	CreateBy *string `json:"createBy"` // optional
	Customer struct {
		ExternalId *string `json:"externalId"` // optional
		Email      *string `json:"email"`      // optional
		Phone      *string `json:"phone"`      // optional
	} `json:"customer"`
}

// response
type CustomerTIResp struct {
	Create     bool   `json:"create"`
	CustomerId string `json:"customerId"`
	Jwt        string `json:"jwt"`
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

func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func AuthenticatePat(ctx context.Context, db *pgxpool.Pool, token string) (Account, error) {
	var account Account

	stmt := `SELECT 
		a.account_id, a.email,
		a.provider, a.auth_user_id, a.name,
		a.created_at, a.updated_at
		FROM account a
		INNER JOIN account_pat ap ON a.account_id = ap.account_id
		WHERE ap.token = $1`

	row, err := db.Query(ctx, stmt, token)
	if err != nil {
		return account, err
	}
	defer row.Close()

	if !row.Next() {
		fmt.Printf("no linked account found for token: %s\n", token)
		return account, sql.ErrNoRows
	}

	err = row.Scan(
		&account.AccountId, &account.Email,
		&account.Provider, &account.AuthUserId, &account.Name,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan linked account for token: %s with error: %v\n", token, err)
		return account, err
	}

	return account, nil
}

func (a Account) GenId() string {
	return "a_" + xid.New().String()
}

func (a Account) GetOrCreateByAuthUserId(ctx context.Context, db *pgxpool.Pool) (Account, bool, error) {
	aId := a.GenId()

	st := `WITH ins AS (
		INSERT INTO account(account_id, auth_user_id, email, provider, name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (auth_user_id) DO NOTHING
		RETURNING
		account_id, auth_user_id, email,
		provider, name,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT account_id, auth_user_id, email, provider, name,
	created_at, updated_at, FALSE AS is_created FROM account
	WHERE auth_user_id = $2 AND NOT EXISTS (SELECT 1 FROM ins)`

	var account Account
	var isCreated bool
	row, err := db.Query(ctx, st, aId, a.AuthUserId, a.Email, a.Provider, a.Name)
	if err != nil {
		return account, isCreated, err
	}
	defer row.Close()

	if !row.Next() {
		return account, isCreated, sql.ErrNoRows
	}

	err = row.Scan(
		&account.AccountId, &account.AuthUserId,
		&account.Email, &account.Provider, &account.Name,
		&account.CreatedAt, &account.UpdatedAt,
		&isCreated,
	)
	if err != nil {
		return account, isCreated, err
	}

	return account, isCreated, nil

}

func (a Account) GetByAuthUserId(ctx context.Context, db *pgxpool.Pool) (Account, error) {
	var account Account

	row, err := db.Query(ctx, `SELECT 
		account_id, auth_user_id, email,
		provider, name, created_at, updated_at
		FROM account WHERE auth_user_id = $1`, a.AuthUserId)
	if err != nil {
		return account, err
	}
	defer row.Close()

	if !row.Next() {
		return account, sql.ErrNoRows
	}

	err = row.Scan(
		&account.AccountId, &account.AuthUserId,
		&account.Email, &account.Provider,
		&account.Name, &account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (c Customer) GenId() string {
	return "c_" + xid.New().String()
}

func (c Customer) GetByExtId(ctx context.Context, db *pgxpool.Pool, workspaceId string, extId string) (Customer, error) {
	var customer Customer

	row, err := db.Query(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND external_id = $2`, workspaceId, extId)
	if err != nil {
		return customer, err
	}
	defer row.Close()

	if !row.Next() {
		return customer, sql.ErrNoRows
	}

	err = row.Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)
	if err != nil {
		return customer, err
	}

	return customer, nil
}

func (c Customer) GetByEmail(ctx context.Context, db *pgxpool.Pool, workspaceId string, email string) (Customer, error) {
	var customer Customer

	row, err := db.Query(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND email = $2`, workspaceId, email)
	if err != nil {
		return customer, err
	}
	defer row.Close()

	if !row.Next() {
		return customer, sql.ErrNoRows
	}

	err = row.Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)
	if err != nil {
		return customer, err
	}

	return customer, nil
}

func (c Customer) GetByPhone(ctx context.Context, db *pgxpool.Pool, workspaceId string, phone string) (Customer, error) {
	var customer Customer

	row, err := db.Query(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND phone = $2`, workspaceId, phone)
	if err != nil {
		return customer, err
	}
	defer row.Close()

	if !row.Next() {
		return customer, sql.ErrNoRows
	}

	err = row.Scan(
		&customer.WorkspaceId, &customer.CustomerId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.Name,
		&customer.CreatedAt, &customer.UpdatedAt,
	)
	if err != nil {
		return customer, err
	}

	return customer, nil
}

func (c Customer) GetOrCreateByExtId(ctx context.Context, db *pgxpool.Pool) (Customer, bool, error) {

	cId := c.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (workspace_id, external_id) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone,
	created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, external_id) = ($2, $3) AND NOT EXISTS (SELECT 1 FROM ins)`

	var customer Customer
	var isCreated bool
	row, err := db.Query(ctx, st, cId, c.WorkspaceId, c.ExternalId, c.Email, c.Phone)
	if err != nil {
		return customer, isCreated, err
	}
	defer row.Close()

	if !row.Next() {
		return customer, isCreated, sql.ErrNoRows
	}

	err = row.Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)
	if err != nil {
		return customer, isCreated, err
	}

	return customer, isCreated, nil
}

func (c Customer) GetOrCreateByEmail(ctx context.Context, db *pgxpool.Pool) (Customer, bool, error) {

	cId := c.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (workspace_id, email) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone,
	created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, email) = ($2, $4) AND NOT EXISTS (SELECT 1 FROM ins)`

	var customer Customer
	var isCreated bool
	row, err := db.Query(ctx, st, cId, c.WorkspaceId, c.ExternalId, c.Email, c.Phone)
	if err != nil {
		return customer, isCreated, err
	}
	defer row.Close()

	if !row.Next() {
		return customer, isCreated, sql.ErrNoRows
	}

	err = row.Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)
	if err != nil {
		return customer, isCreated, err
	}

	return customer, isCreated, nil
}

func (c Customer) GetOrCreateByPhone(ctx context.Context, db *pgxpool.Pool) (Customer, bool, error) {

	cId := c.GenId()
	st := `WITH ins AS (
		INSERT INTO customer (customer_id, workspace_id, external_id, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (workspace_id, phone) DO NOTHING
		RETURNING
		customer_id, workspace_id,
		external_id, email, phone,
		created_at, updated_at,
		TRUE AS is_created
	)
	SELECT * FROM ins
	UNION ALL
	SELECT customer_id, workspace_id, external_id, email, phone,
	created_at, updated_at, FALSE AS is_created FROM customer
	WHERE (workspace_id, phone) = ($2, $5) AND NOT EXISTS (SELECT 1 FROM ins)`

	var customer Customer
	var isCreated bool
	row, err := db.Query(ctx, st, cId, c.WorkspaceId, c.ExternalId, c.Email, c.Phone)
	if err != nil {
		return customer, isCreated, err
	}
	defer row.Close()

	if !row.Next() {
		return customer, isCreated, sql.ErrNoRows
	}

	err = row.Scan(
		&customer.CustomerId, &customer.WorkspaceId,
		&customer.ExternalId, &customer.Email,
		&customer.Phone, &customer.CreatedAt,
		&customer.UpdatedAt, &isCreated,
	)
	if err != nil {
		return customer, isCreated, err
	}

	return customer, isCreated, nil
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func AuthenticateMember(ctx context.Context, db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) (Account, error) {
	var account Account

	ath := r.Header.Get("Authorization")
	if ath == "" {
		return account, fmt.Errorf("cannot authenticate without authorization header")
	}

	cred := strings.Split(ath, " ")
	scheme := strings.ToLower(cred[0])

	if scheme == "token" {
		fmt.Println("authenticate via PAT")
		account, err := AuthenticatePat(ctx, db, cred[1])
		if err != nil {
			return account, fmt.Errorf("failed to authenticate with error: %v", err)
		}
		fmt.Printf("authenticated account with accountId: %s\n", account.AccountId)
		return account, nil
	} else if scheme == "bearer" {
		fmt.Println("authenticate via JWTs")
		hmacSecret, err := GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			return account, fmt.Errorf("failed to get env SUPABASE_JWT_SECRET with error: %v", err)
		}
		auid, err := parseJWTToken(cred[1], []byte(hmacSecret))
		if err != nil {
			return account, fmt.Errorf("failed to parse JWT token with error: %v", err)
		}
		fmt.Printf("authenticated account with id: %s\n", auid.Id)
		ta := Account{AuthUserId: auid.Id, Email: auid.Email, Provider: "supabase"}
		account, err := ta.GetByAuthUserId(ctx, db)
		if err != nil {
			return account, fmt.Errorf("failed to get account by authUserId: %s with error: %v", ta.AuthUserId, err)
		}
		return account, nil
	} else {
		return account, fmt.Errorf("unsupported scheme: `%s` cannot authenticate", scheme)
	}
}

// JWTClaims taken from Supabase JWT encoding
type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type AuthUserId struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func parseJWTToken(token string, hmacSecret []byte) (auid AuthUserId, err error) {
	t, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return auid, fmt.Errorf("error validating jwt token with error: %v", err)
	} else if claims, ok := t.Claims.(*JWTClaims); ok {
		sub, err := claims.RegisteredClaims.GetSubject()
		if err != nil {
			return auid, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}
		return AuthUserId{Id: sub, Email: claims.Email}, nil
	}

	return auid, fmt.Errorf("error parsing token: %v", token)
}

func handleAuthAccountMaker(ctx context.Context, db *pgxpool.Pool) http.Handler {
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
			hmacSecret, err := GetEnv("SUPABASE_JWT_SECRET")
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			auid, err := parseJWTToken(cred[1], []byte(hmacSecret))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ta := Account{AuthUserId: auid.Id, Email: auid.Email, Provider: "supabase"}
			account, isCreated, err := ta.GetOrCreateByAuthUserId(ctx, db)
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

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleGetWorkspaces(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		account, err := AuthenticateMember(ctx, db, w, r)
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

		workspaces := make([]Workspace, 0)
		for rows.Next() {
			var workspace Workspace
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

func handleLLMQuery(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var workspace Workspace

		var query LLMRequestQuery
		err := json.NewDecoder(r.Body).Decode(&query)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		workspaceId := r.PathValue("workspaceId")

		row, err := db.Query(ctx, `SELECT workspace_id, account_id,
			name, created_at, updated_at
			FROM workspace WHERE workspace_id = $1`,
			workspaceId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer row.Close()

		if !row.Next() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		err = row.Scan(
			&workspace.WorkspaceId, &workspace.AccountId,
			&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		requestId := xid.New()
		wrkLLM := LLM{WorkspaceId: workspaceId, Prompt: query.Q, RequestId: requestId.String()}

		resp, err := wrkLLM.Generate()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(ctx, `INSERT INTO
			llm_rr_log(workspace_id, request_id, prompt, response, model)
			VALUES ($1, $2, $3, $4, $5)`,
			workspace.WorkspaceId, wrkLLM.RequestId, wrkLLM.Prompt, resp.Text, resp.Model)
		if err != nil {
			log.Printf("failed to insert into llm request response log with error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})
}

func handleLLMQueryEval(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var workspace Workspace

		var eval LLMRREval
		err := json.NewDecoder(r.Body).Decode(&eval)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		workspaceId := r.PathValue("workspaceId")

		row, err := db.Query(ctx, `SELECT workspace_id, account_id,
			name, created_at, updated_at
			FROM workspace WHERE workspace_id = $1`,
			workspaceId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer row.Close()

		if !row.Next() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		err = row.Scan(
			&workspace.WorkspaceId, &workspace.AccountId,
			&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		requestId := r.PathValue("requestId")

		_, err = db.Exec(ctx, `UPDATE llm_rr_log SET eval = $1
			WHERE workspace_id = $2 AND request_id = $3`,
			eval.Eval, workspace.WorkspaceId, requestId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})
}

func handleCustomerTokenIssue(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		worskpaceId := r.PathValue("workspaceId")
		fmt.Printf("issue token for customer in workspaceId: %v\n", worskpaceId)

		var rb CustomerTIRequest
		err := json.NewDecoder(r.Body).Decode(&rb)
		if err != nil {
			fmt.Printf("failed to decode request body error: %v\n", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		tc := Customer{
			WorkspaceId: worskpaceId,
		}
		tc.ExternalId = NullString(rb.Customer.ExternalId)
		tc.Email = NullString(rb.Customer.Email)
		tc.Phone = NullString(rb.Customer.Phone)

		if !tc.ExternalId.Valid && !tc.Email.Valid && !tc.Phone.Valid {
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
				if !tc.Email.Valid {
					fmt.Println("`email` is required for `createBy` email")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				email := tc.Email.String
				fmt.Printf("create the customer by email %s\n", email)
				customer, isCreated, err := tc.GetOrCreateByEmail(ctx, db)
				if err != nil {
					fmt.Printf("failed to get or create customer by email %s with error: %v\n", email, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				fmt.Printf("customerId: %s is created: %v\n", customer.CustomerId, isCreated)
				resp := CustomerTIResp{
					Create:     isCreated,
					CustomerId: customer.CustomerId,
					Jwt:        "TODO: generate jwt token",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				return
			case "phone":
				if !tc.Phone.Valid {
					fmt.Println("`phone` is required for `createBy` phone")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				phone := tc.Phone.String
				fmt.Printf("create the customer by phone %s\n", phone)
				customer, isCreated, err := tc.GetOrCreateByPhone(ctx, db)
				if err != nil {
					fmt.Printf("failed to get or create customer by phone %s with error: %v\n", phone, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				fmt.Printf("customerId: %s is created: %v\n", customer.CustomerId, isCreated)
				resp := CustomerTIResp{
					Create:     isCreated,
					CustomerId: customer.CustomerId,
					Jwt:        "TODO: generate jwt token",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			case "externalId":
				if !tc.ExternalId.Valid {
					fmt.Println("`externalId` is required for `createBy` externalId")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				extId := tc.ExternalId.String
				fmt.Printf("create the customer by externalId %s\n", extId)
				customer, isCreated, err := tc.GetOrCreateByExtId(ctx, db)
				if err != nil {
					fmt.Printf("failed to get or create customer by externalId %s with error: %v\n", extId, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				fmt.Printf("customerId: %s is created: %v\n", customer.CustomerId, isCreated)
				resp := CustomerTIResp{
					Create:     isCreated,
					CustomerId: customer.CustomerId,
					Jwt:        "TODO: generate jwt token",
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
			var customer Customer
			fmt.Printf("based on identifiers check for customer in workspaceId: %v\n", worskpaceId)
			if tc.ExternalId.Valid {
				extId := tc.ExternalId.String
				fmt.Printf("get customer by externalId %s\n", extId)
				customer, err = customer.GetByExtId(ctx, db, worskpaceId, extId)
				if err != nil {
					fmt.Printf("failed to get customer by externalId %s with error: %v\n", extId, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				fmt.Printf("found customer with customer id: %s\n", customer.CustomerId)
				resp := CustomerTIResp{
					Create:     false,
					CustomerId: customer.CustomerId,
					Jwt:        "TODO: generate jwt token",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else if tc.Email.Valid {
				email := tc.Email.String
				fmt.Printf("get customer by email %s\n", email)
				customer, err = customer.GetByEmail(ctx, db, worskpaceId, email)
				if err != nil {
					fmt.Printf("failed to get customer by email %s with error: %v\n", email, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				fmt.Printf("found customer with customer id: %s\n", customer.CustomerId)
				resp := CustomerTIResp{
					Create:     false,
					CustomerId: customer.CustomerId,
					Jwt:        "TODO: generate jwt token",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else if tc.Phone.Valid {
				phone := tc.Phone.String
				fmt.Printf("get customer by phone %s\n", phone)
				customer, err = customer.GetByPhone(ctx, db, worskpaceId, phone)
				if err != nil {
					fmt.Printf("failed to get customer by phone %s with error: %v\n", phone, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				fmt.Printf("found customer with customer id: %s\n", customer.CustomerId)
				resp := CustomerTIResp{
					Create:     false,
					CustomerId: customer.CustomerId,
					Jwt:        "TODO: generate jwt token",
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

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var err error

	pgConnStr, err := GetEnv("POSTGRES_URI")
	if err != nil {
		return err
	}

	db, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		return fmt.Errorf("unable to create pg connection pool: %v", err)
	}

	defer db.Close()

	var tm time.Time

	err = db.QueryRow(ctx, "SELECT NOW()").Scan(&tm)

	if err != nil {
		return fmt.Errorf("failed to query database: %v", err)
	}

	log.Printf("database ready with db time: %s\n", tm.Format(time.RFC1123))

	mux := http.NewServeMux()

	// sdk+web
	mux.Handle("POST /accounts/auth/{$}", handleAuthAccountMaker(ctx, db))

	// sdk+web - TODO: add member authentication
	mux.HandleFunc("GET /{$}", handleGetIndex)

	// sdk+web - TODO: add member authentication
	mux.Handle("GET /workspaces/{$}",
		handleGetWorkspaces(ctx, db))

	// customer
	mux.Handle("POST /workspaces/{workspaceId}/-/queries/{$}",
		handleLLMQuery(ctx, db))

	// customer
	mux.Handle("POST /workspaces/{workspaceId}/-/queries/{requestId}/{$}",
		handleLLMQueryEval(ctx, db))

	// sdk+web
	mux.Handle("POST /workspaces/{workspaceId}/tokens/{$}",
		handleCustomerTokenIssue(ctx, db))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders: []string{"*"},
	})

	srv := &http.Server{
		Addr:              *addr,
		Handler:           LoggingMiddleware(c.Handler(mux)),
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
	}

	log.Printf("server up and running on %s", *addr)

	err = srv.ListenAndServe()

	return err
}

func main() {
	flag.Parse()
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
