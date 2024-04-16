package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
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

const DefaultAuthProvider string = "supabase"

var addr = flag.String("addr", "127.0.0.1:8080", "listen address")

func GetEnv(key string) (string, error) {
	value, status := os.LookupEnv(key)
	if !status {
		return "", fmt.Errorf("env `%s` is not set", key)
	}
	return value, nil
}

func GenToken(length int, prefix string) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return prefix + base64.URLEncoding.EncodeToString(buffer)[:length], nil
}

func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// model
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

func (a Account) GenId() string {
	return "a_" + xid.New().String()
}

func (a Account) GetOrCreateByAuthUserId(
	ctx context.Context, db *pgxpool.Pool,
	authUserId string, email string, provider string, name string,
) (Account, bool, error) {

	var account Account
	var isCreated bool

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

	row, err := db.Query(ctx, st, aId, authUserId, email, provider, name)
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

func (a Account) GetByAuthUserId(ctx context.Context, db *pgxpool.Pool, authUserId string) (Account, error) {
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

// model
type AccountPAT struct {
	AccountId   string
	PatId       string
	Token       string
	Name        string
	Description sql.NullString
	UnMask      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ap AccountPAT) MarshalJSON() ([]byte, error) {
	var description *string
	var token string
	if ap.Description.Valid {
		description = &ap.Description.String
	}

	maskLeft := func(s string) string {
		rs := []rune(s)
		for i := range rs[:len(rs)-4] {
			rs[i] = '*'
		}
		return string(rs)
	}

	if !ap.UnMask {
		token = maskLeft(ap.Token)
	} else {
		token = ap.Token
	}

	aux := &struct {
		AccountId   string  `json:"accountId"`
		PatId       string  `json:"patId"`
		Token       string  `json:"token"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}{
		AccountId:   ap.AccountId,
		PatId:       ap.PatId,
		Token:       token,
		Name:        ap.Name,
		Description: description,
		CreatedAt:   ap.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   ap.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (ap AccountPAT) GenId() string {
	return "ap_" + xid.New().String()
}

func (ap AccountPAT) Create(ctx context.Context, db *pgxpool.Pool) (AccountPAT, error) {
	apId := ap.GenId()
	token, err := GenToken(32, "pt_")
	if err != nil {
		return ap, err
	}
	stmt := `INSERT INTO account_pat(account_id, pat_id, token, name, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING account_id, pat_id, token, name, description, created_at, updated_at`

	row, err := db.Query(ctx, stmt, ap.AccountId, apId, token, ap.Name, ap.Description)
	if err != nil {
		return ap, err
	}
	defer row.Close()

	if !row.Next() {
		return ap, sql.ErrNoRows
	}

	err = row.Scan(
		&ap.AccountId, &ap.PatId, &ap.Token,
		&ap.Name, &ap.Description, &ap.CreatedAt, &ap.UpdatedAt,
	)
	if err != nil {
		return ap, err
	}
	return ap, nil
}

func (ap AccountPAT) GetListByAccountId(ctx context.Context, db *pgxpool.Pool) ([]AccountPAT, error) {
	stmt := `SELECT account_id, pat_id, token, name, description,
		created_at, updated_at
		FROM account_pat WHERE account_id = $1
		ORDER BY created_at DESC LIMIT 100`

	rows, err := db.Query(ctx, stmt, ap.AccountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aps := make([]AccountPAT, 0)
	for rows.Next() {
		var ap AccountPAT
		err = rows.Scan(
			&ap.AccountId, &ap.PatId, &ap.Token,
			&ap.Name, &ap.Description, &ap.CreatedAt, &ap.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		aps = append(aps, ap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return aps, nil
}

// model
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

func (w Workspace) GenId() string {
	return "wrk" + xid.New().String()
}

func (w Workspace) Create(ctx context.Context, db *pgxpool.Pool) (Workspace, error) {
	var workspace Workspace
	var member Member

	tx, err := db.Begin(ctx)
	if err != nil {
		fmt.Printf("failed to start db transaction with error: %v\n", err)
		return workspace, err
	}
	defer tx.Rollback(ctx)

	wId := w.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO workspace(workspace_id, account_id, name)
		VALUES ($1, $2, $3)
		RETURNING
		workspace_id, account_id, name, created_at, updated_at`, wId, w.AccountId, w.Name).Scan(
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to insert workspace with error: %v\n", err)
		return workspace, err
	}

	mId := member.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO member(workspace_id, account_id, member_id, role)
		VALUES ($1, $2, $3, $4)
		RETURNING
		workspace_id, account_id, member_id, name, role, created_at, updated_at`,
		wId, workspace.AccountId, mId, "primary").Scan(
		&member.WorkspaceId, &member.AccountId,
		&member.MemberId, &member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to insert member with error: %v\n", err)
		return workspace, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit db transaction with error: %v\n", err)
		return workspace, err
	}
	return workspace, nil
}

func (w Workspace) GetAccountWorkspace(ctx context.Context, db *pgxpool.Pool) (Workspace, error) {
	var workspace Workspace

	row, err := db.Query(ctx, `SELECT 
		workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE account_id = $1 AND workspace_id = $2`, w.AccountId, w.WorkspaceId)
	if err != nil {
		return workspace, err
	}
	defer row.Close()

	if !row.Next() {
		return workspace, sql.ErrNoRows
	}

	err = row.Scan(
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)
	if err != nil {
		return workspace, err
	}

	return workspace, nil
}

func (w Workspace) GetById(ctx context.Context, db *pgxpool.Pool) (Workspace, error) {
	var workspace Workspace

	row, err := db.Query(ctx, `SELECT 
		workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE workspace_id = $1`, w.WorkspaceId)
	if err != nil {
		return workspace, err
	}
	defer row.Close()

	if !row.Next() {
		return workspace, sql.ErrNoRows
	}

	err = row.Scan(
		&workspace.WorkspaceId, &workspace.AccountId,
		&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt,
	)
	if err != nil {
		return workspace, err
	}

	return workspace, nil
}

// model
type Member struct {
	WorkspaceId string
	AccountId   string
	MemberId    string
	Name        sql.NullString
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m Member) GenId() string {
	return "m_" + xid.New().String()
}

func (m Member) MarshalJSON() ([]byte, error) {
	var name *string
	if m.Name.Valid {
		name = &m.Name.String
	}
	aux := &struct {
		WorkspaceId string  `json:"workspaceId"`
		AccountId   string  `json:"accountId"`
		MemberId    string  `json:"memberId"`
		Name        *string `json:"name"`
		Role        string  `json:"role"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}{
		WorkspaceId: m.WorkspaceId,
		AccountId:   m.AccountId,
		MemberId:    m.MemberId,
		Name:        name,
		Role:        m.Role,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (m Member) GetWorkspaceMemberByAccountId(
	ctx context.Context, db *pgxpool.Pool, workspaceId string, accountId string,
) (Member, error) {

	row, err := db.Query(ctx, `SELECT
		workspace_id, account_id, member_id, name, role, created_at, updated_at
		FROM member WHERE workspace_id = $1 AND account_id = $2`, workspaceId, accountId)
	if err != nil {
		return m, err
	}

	defer row.Close()

	if !row.Next() {
		return m, sql.ErrNoRows
	}

	err = row.Scan(
		&m.WorkspaceId, &m.AccountId,
		&m.MemberId, &m.Name, &m.Role,
		&m.CreatedAt, &m.UpdatedAt,
	)

	if err != nil {
		return m, err
	}

	return m, nil
}

// model
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

func (c Customer) GenId() string {
	return "c_" + xid.New().String()
}

func (c Customer) GetById(ctx context.Context, db *pgxpool.Pool) (Customer, error) {
	var customer Customer

	row, err := db.Query(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND customer_id = $2`, c.WorkspaceId, c.CustomerId)
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

func (c Customer) MakeJWT() (string, error) {

	var externalId string
	var email string
	var phone string

	audience := []string{"customer"}

	sk, err := GetEnv("SUPABASE_JWT_SECRET")
	if err != nil {
		return "", fmt.Errorf("failed to get env SUPABASE_JWT_SECRET with error: %v", err)
	}

	if !c.ExternalId.Valid {
		externalId = ""
	} else {
		externalId = c.ExternalId.String
	}

	if !c.Email.Valid {
		email = ""
	} else {
		email = c.Email.String
	}

	if !c.Phone.Valid {
		phone = ""
	} else {
		phone = c.Phone.String
	}

	claims := CustomerJWTClaims{
		WorkspaceId: c.WorkspaceId,
		ExternalId:  externalId,
		Email:       email,
		Phone:       phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.zyg.ai",
			Subject:   c.CustomerId,
			Audience:  audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(1, 0, 0)), // Set ExpiresAt to 1 year from now
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        c.WorkspaceId + ":" + c.CustomerId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString([]byte(sk))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token with error: %v", err)
	}
	return jwt, nil
}

// model
type ThreadQA struct {
	WorkspaceId    string
	CustomerId     string
	ThreadId       string
	ParentThreadId sql.NullString
	Query          string
	Title          string
	Summary        string
	Sequence       int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (tq ThreadQA) MarshalJSON() ([]byte, error) {
	var pth *string
	if tq.ParentThreadId.Valid {
		pth = &tq.ParentThreadId.String
	}
	aux := &struct {
		WorkspaceId    string  `json:"workspaceId"`
		CustomerId     string  `json:"customerId"`
		ThreadId       string  `json:"threadId"`
		ParentThreadId *string `json:"parentThreadId"`
		Query          string  `json:"query"`
		Title          string  `json:"title"`
		Summary        string  `json:"summary"`
		Sequence       int     `json:"sequence"`
		CreatedAt      string  `json:"createdAt"`
		UpdatedAt      string  `json:"updatedAt"`
	}{
		WorkspaceId:    tq.WorkspaceId,
		CustomerId:     tq.CustomerId,
		ThreadId:       tq.ThreadId,
		ParentThreadId: pth,
		Query:          tq.Query,
		Title:          tq.Title,
		Summary:        tq.Summary,
		Sequence:       tq.Sequence,
		CreatedAt:      tq.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      tq.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (tq ThreadQA) GenId() string {
	return "tq_" + xid.New().String()
}

func (tq ThreadQA) Create(ctx context.Context, db *pgxpool.Pool) (ThreadQA, error) {
	var thread ThreadQA

	tqId := tq.GenId()
	stmt := `INSERT INTO 
		thread_qa(workspace_id, customer_id, thread_id, parent_thread_id, query, title, summary)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING 
		workspace_id, customer_id, thread_id, parent_thread_id,
		query, title, summary, sequence,
		created_at, updated_at`

	row, err := db.Query(ctx, stmt, tq.WorkspaceId, tq.CustomerId, tqId, tq.ParentThreadId, tq.Query, tq.Title, tq.Summary)
	if err != nil {
		return thread, err
	}
	defer row.Close()

	if !row.Next() {
		return thread, sql.ErrNoRows
	}

	err = row.Scan(
		&thread.WorkspaceId, &thread.CustomerId, &thread.ThreadId, &thread.ParentThreadId,
		&thread.Query, &thread.Title, &thread.Summary, &thread.Sequence,
		&thread.CreatedAt, &thread.UpdatedAt,
	)
	if err != nil {
		return thread, err
	}

	return thread, nil
}

// model
type ThreadQAA struct {
	WorkspaceId string
	ThreadQAId  string
	AnswerId    string
	Answer      string
	Sequence    int
	Eval        sql.NullInt32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (tqa ThreadQAA) MarshalJSON() ([]byte, error) {
	var eval *int32
	if tqa.Eval.Valid {
		eval = &tqa.Eval.Int32
	}
	aux := &struct {
		WorkspaceId string `json:"workspaceId"`
		ThreadQAId  string `json:"threadQAId"`
		AnswerId    string `json:"answerId"`
		Answer      string `json:"answer"`
		Sequence    int    `json:"sequence"`
		Eval        *int32 `json:"eval"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}{
		WorkspaceId: tqa.WorkspaceId,
		ThreadQAId:  tqa.ThreadQAId,
		AnswerId:    tqa.AnswerId,
		Answer:      tqa.Answer,
		Sequence:    tqa.Sequence,
		Eval:        eval,
		CreatedAt:   tqa.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   tqa.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}
func (tqa ThreadQAA) GenId() string {
	return "tqa_" + xid.New().String()
}

func (tqa ThreadQAA) Create(ctx context.Context, db *pgxpool.Pool) (ThreadQAA, error) {
	var thread ThreadQAA

	tqaId := tqa.GenId()
	stmt := `INSERT INTO 
		thread_qa_answer(workspace_id, thread_qa_id, answer_id, answer)
		VALUES ($1, $2, $3, $4)
		RETURNING 
		workspace_id, thread_qa_id, answer_id, answer, 
		eval, sequence, created_at, updated_at`

	row, err := db.Query(ctx, stmt, tqa.WorkspaceId, tqa.ThreadQAId, tqaId, tqa.Answer)
	if err != nil {
		return thread, err
	}
	defer row.Close()

	if !row.Next() {
		return thread, sql.ErrNoRows
	}

	err = row.Scan(
		&thread.WorkspaceId, &thread.ThreadQAId, &thread.AnswerId, &thread.Answer,
		&thread.Eval, &thread.Sequence, &thread.CreatedAt, &thread.UpdatedAt,
	)
	if err != nil {
		return thread, err
	}

	return thread, nil
}

// model
type ThreadChat struct {
	WorkspaceId  string
	CustomerId   string
	CustomerName sql.NullString
	AssigneeId   sql.NullString
	AssigneeName sql.NullString
	ThreadChatId string
	Title        string
	Summary      string
	Sequence     int
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (th ThreadChat) GenId() string {
	return "th_" + xid.New().String()
}

func (th ThreadChat) CreateCustomerThChat(
	ctx context.Context, db *pgxpool.Pool, w Workspace, c Customer, m string,
) (ThreadChat, ThreadChatMessage, error) {
	var thm ThreadChatMessage

	tx, err := db.Begin(ctx)
	if err != nil {
		fmt.Printf("failed to start db transaction with error: %v\n", err)
		return th, thm, err
	}

	defer tx.Rollback(ctx)

	thId := th.GenId()
	th.Status = "todo" // default status
	err = tx.QueryRow(ctx, `INSERT INTO thread_chat(workspace_id, customer_id, thread_chat_id, title, summary, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
		workspace_id, customer_id, assignee_id,
		thread_chat_id, title, summary, sequence, status, created_at, updated_at`,
		w.WorkspaceId, c.CustomerId, thId, th.Title, th.Summary, th.Status).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.AssigneeId,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to insert in thread chat with error: %v\n", err)
		return th, thm, err
	}

	// thm.ThreadChatMessageId = thm.GenId()
	thmId := thm.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO thread_chat_message(thread_chat_id, thread_chat_message_id, body, customer_id)
		VALUES ($1, $2, $3, $4)
		RETURNING
		thread_chat_id, thread_chat_message_id, body, sequence, customer_id, member_id, created_at, updated_at`,
		th.ThreadChatId, thmId, m, c.CustomerId).Scan(
		&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
		&thm.Sequence, &thm.CustomerId, &thm.MemberId,
		&thm.CreatedAt, &thm.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to insert thread chat message with error: %v\n", err)
		return th, thm, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit db transaction with error: %v\n", err)
		return th, thm, err
	}

	// set customer attributes we already have
	th.CustomerName = c.Name
	thm.CustomerName = c.Name

	return th, thm, nil
}

func (th ThreadChat) GetById(ctx context.Context, db *pgxpool.Pool, threadId string) (ThreadChat, error) {
	stmt := `SELECT th.workspace_id AS workspace_id,
		c.customer_id AS customer_id,
		c.name AS customer_name,
		a.member_id AS assignee_id,
		a.name AS assignee_name,
		th.thread_chat_id AS thread_chat_id,
		th.title AS title,
		th.summary AS summary,
		th.sequence AS sequence,
		th.status AS status,
		th.created_at AS created_at,
		th.updated_at AS updated_at
		FROM thread_chat th
		INNER JOIN customer c ON th.customer_id = c.customer_id
		LEFT OUTER JOIN member a ON th.assignee_id = a.member_id
		WHERE th.thread_chat_id = $1`

	err := db.QueryRow(ctx, stmt, threadId).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
	)
	if err != nil {
		return th, err
	}

	return th, nil
}

func (th ThreadChat) GetListByWorkspaceCustomerId(
	ctx context.Context, db *pgxpool.Pool, workspaceId string, customerId string,
) ([]ThreadChatWithMessage, error) {

	ths := make([]ThreadChatWithMessage, 0, 100)

	stmt := `SELECT
			th.workspace_id AS workspace_id,	
			thc.customer_id AS thread_customer_id,
			thc.name AS thread_customer_name,
			tha.member_id AS thread_assignee_id,
			tha.name AS thread_assignee_name,
			th.thread_chat_id AS thread_chat_id,
			th.title AS title,
			th.summary AS summary,
			th.sequence AS sequence,
			th.status AS status,
			th.created_at AS created_at,
			th.updated_at AS updated_at,
			thm.thread_chat_id AS message_thread_chat_id,
			thm.thread_chat_message_id AS thread_chat_message_id,
			thm.body AS message_body,
			thm.sequence AS message_sequence,
			thm.created_at AS message_created_at,
			thm.updated_at AS message_updated_at,
			thmc.customer_id AS message_customer_id,
			thmc.name AS message_customer_name,
			thmm.member_id AS message_member_id,
			thmm.name AS message_member_name
		FROM thread_chat th
		INNER JOIN thread_chat_message thm ON th.thread_chat_id = thm.thread_chat_id
		INNER JOIN customer thc ON th.customer_id = thc.customer_id
		LEFT OUTER JOIN member tha ON th.assignee_id = tha.member_id
		LEFT OUTER JOIN customer thmc ON thm.customer_id = thmc.customer_id
		LEFT OUTER JOIN member thmm ON thm.member_id = thmm.member_id
		INNER JOIN (
			SELECT thread_chat_id, MAX(sequence) AS sequence
			FROM thread_chat_message
			GROUP BY
			thread_chat_id
		) latest ON thm.thread_chat_id = latest.thread_chat_id
		AND thm.sequence = latest.sequence
		WHERE th.workspace_id = $1 AND th.customer_id = $2
		ORDER BY sequence DESC LIMIT 100`

	rows, err := db.Query(ctx, stmt, workspaceId, customerId)
	if err != nil {
		return ths, err
	}

	defer rows.Close()

	for rows.Next() {
		var th ThreadChat
		var tc ThreadChatMessage
		err = rows.Scan(
			&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
			&th.AssigneeId, &th.AssigneeName,
			&th.ThreadChatId, &th.Title, &th.Summary,
			&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
			&tc.ThreadChatId, &tc.ThreadChatMessageId, &tc.Body,
			&tc.Sequence, &tc.CreatedAt, &tc.UpdatedAt,
			&tc.CustomerId, &tc.CustomerName, &tc.MemberId, &tc.MemberName,
		)
		if err != nil {
			return ths, err
		}

		thm := ThreadChatWithMessage{
			ThreadChat: th,
			Message:    tc,
		}
		ths = append(ths, thm)
	}

	if err := rows.Err(); err != nil {
		return ths, err
	}

	return ths, nil
}

func (th ThreadChat) AssignMember(
	ctx context.Context, db *pgxpool.Pool, memberId string, threadChatId string,
) (ThreadChat, error) {
	stmt := `WITH ups AS (
		UPDATE thread_chat SET assignee_id = $1
		WHERE thread_chat_id = $2
		RETURNING
		workspace_id, thread_chat_id, customer_id, assignee_id,
		title, summary, sequence, status, created_at, updated_at
	) SELECT
	ups.workspace_id AS workspace_id,
	c.customer_id AS customer_id,
	c.name AS customer_name,
	m.member_id AS assignee_id,
	m.name AS assignee_name,
	ups.thread_chat_id AS thread_chat_id,
	ups.title AS title,
	ups.summary AS summary,
	ups.sequence AS sequence,
	ups.status AS status,
	ups.created_at AS created_at,
	ups.updated_at AS updated_at
	FROM ups
	INNER JOIN customer c ON ups.customer_id = c.customer_id
	LEFT OUTER JOIN member m ON ups.assignee_id = m.member_id`

	row, err := db.Query(ctx, stmt, memberId, threadChatId)
	if err != nil {
		return th, err
	}

	defer row.Close()

	if !row.Next() {
		return th, sql.ErrNoRows
	}

	err = row.Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
	)

	if err != nil {
		return th, err
	}

	return th, nil
}

// model
type ThreadChatMessage struct {
	ThreadChatId        string
	ThreadChatMessageId string
	Body                string
	Sequence            int
	CustomerId          sql.NullString
	CustomerName        sql.NullString
	MemberId            sql.NullString
	MemberName          sql.NullString
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (thm ThreadChatMessage) GenId() string {
	return "thm_" + xid.New().String()
}

func (thm ThreadChatMessage) CreateCustomerThChatMessage(
	ctx context.Context, db *pgxpool.Pool, th ThreadChat, c Customer, msg string,
) (ThreadChatMessage, error) {
	// var customerId sql.NullString

	// if c.CustomerId != "" {
	// 	customerId = sql.NullString{String: c.CustomerId, Valid: true}
	// } else {
	// 	return thm, fmt.Errorf("cannot create customer thread chat message without customer id got: %v", c.CustomerId)
	// }

	// thm.ThreadChatMessageId = thm.GenId()
	// err := db.QueryRow(ctx, `INSERT INTO thread_chat_message(thread_chat_id, thread_chat_message_id, body, customer_id)
	// 	VALUES ($1, $2, $3, $4)
	// 	RETURNING
	// 	thread_chat_id, thread_chat_message_id, body, sequence, customer_id, member_id, created_at, updated_at`,
	// 	th.ThreadChatId, thm.ThreadChatMessageId, msg, customerId).Scan(
	// 	&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
	// 	&thm.Sequence, &thm.CustomerId, &thm.MemberId,
	// 	&thm.CreatedAt, &thm.UpdatedAt,
	// )
	// if err != nil {
	// 	fmt.Printf("failed to insert thread_chat_message with error: %v\n", err)
	// 	return thm, err
	// }
	// // attach customer name before returning
	// // doing this will not require any fancy db query as customer name
	// // is already available and its of the same type.
	// thm.CustomerName = c.Name
	// return thm, nil

	thmId := thm.GenId()
	stmt := `WITH ins AS (
		INSERT INTO thread_chat_message (thread_chat_id, thread_chat_message_id, body, customer_id)
			VALUES ($1, $2, $3, $4)
		RETURNING
			thread_chat_id, thread_chat_message_id, body, sequence,
			customer_id, member_id, created_at, updated_at
		) SELECT ins.thread_chat_id AS thread_chat_id,
		ins.thread_chat_message_id AS thread_chat_message_id,
		ins.body AS body,
		ins.sequence AS sequence,
		c.customer_id AS customer_id,
		c.name AS customer_name,
		m.member_id AS member_id,
		m.name AS member_name,
		ins.created_at AS created_at,
		ins.updated_at AS updated_at
		FROM ins
		LEFT OUTER JOIN customer c ON ins.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ins.member_id = m.member_id`
	err := db.QueryRow(ctx, stmt, th.ThreadChatId, thmId, msg, c.CustomerId).Scan(
		&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
		&thm.Sequence, &thm.CustomerId, &thm.CustomerName,
		&thm.MemberId, &thm.MemberName,
		&thm.CreatedAt, &thm.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed insert for thread chat message for thread id %s with error: %v\n", th.ThreadChatId, err)
		return thm, err
	}
	return thm, nil
}

func (thm ThreadChatMessage) CreateMemberThChatMessage(
	ctx context.Context, db *pgxpool.Pool, th ThreadChat, m Member, msg string,
) (ThreadChatMessage, error) {

	thmId := thm.GenId()
	stmt := `WITH ins AS (
		INSERT INTO thread_chat_message (thread_chat_id, thread_chat_message_id, body, member_id)
			VALUES ($1, $2, $3, $4)
		RETURNING
			thread_chat_id, thread_chat_message_id, body, sequence,
			customer_id, member_id, created_at, updated_at
		) SELECT ins.thread_chat_id AS thread_chat_id,
		ins.thread_chat_message_id AS thread_chat_message_id,
		ins.body AS body,
		ins.sequence AS sequence,
		c.customer_id AS customer_id,
		c.name AS customer_name,
		m.member_id AS member_id,
		m.name AS member_name,
		ins.created_at AS created_at,
		ins.updated_at AS updated_at
		FROM ins
		LEFT OUTER JOIN customer c ON ins.customer_id = c.customer_id
		LEFT OUTER JOIN member m ON ins.member_id = m.member_id`
	err := db.QueryRow(ctx, stmt, th.ThreadChatId, thmId, msg, m.MemberId).Scan(
		&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
		&thm.Sequence, &thm.CustomerId, &thm.CustomerName,
		&thm.MemberId, &thm.MemberName,
		&thm.CreatedAt, &thm.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed insert for thread chat message for thread id %s with error: %v\n", th.ThreadChatId, err)
		return thm, err
	}
	return thm, nil
}

func (thm ThreadChatMessage) GetListByThreadId(
	ctx context.Context, db *pgxpool.Pool, threadId string,
) ([]ThreadChatMessage, error) {

	stmt := `SELECT
		thm.thread_chat_id AS thread_chat_id,
		thm.thread_chat_message_id AS thread_chat_message_id,
		thm.body AS body,
		thm.sequence AS sequence,
		thm.created_at AS created_at,
		thm.updated_at AS updated_at,
		c.customer_id AS customer_id,
		c.name AS customer_name,
		m.member_id AS member_id,
		m.name AS member_name
		FROM thread_chat_message AS thm
		LEFT OUTER JOIN customer AS c ON thm.customer_id = c.customer_id
		LEFT OUTER JOIN member AS m ON thm.member_id = m.member_id
		WHERE thm.thread_chat_id = $1
		ORDER BY sequence DESC LIMIT 100`

	rows, err := db.Query(ctx, stmt, threadId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]ThreadChatMessage, 0, 100)
	for rows.Next() {
		var m ThreadChatMessage
		err = rows.Scan(
			&m.ThreadChatId, &m.ThreadChatMessageId, &m.Body, &m.Sequence,
			&m.CreatedAt, &m.UpdatedAt, &m.CustomerId, &m.CustomerName,
			&m.MemberId, &m.MemberName,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

type ThreadChatWithMessage struct {
	ThreadChat ThreadChat
	Message    ThreadChatMessage
}

// request
type ThreadQAReq struct {
	Query string `json:"query"`
}

// response
type ThreadQAAResp struct {
	AnswerId string
	Answer   string
	Eval     sql.NullInt32
	Sequence int
}

func (tqar ThreadQAAResp) MarshalJSON() ([]byte, error) {
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

type ThreadQAResp struct {
	ThreadId  string
	Query     string
	Sequence  int
	CreatedAt time.Time
	UpdatedAt time.Time
	Answers   []ThreadQAAResp
}

func (thr ThreadQAResp) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ThreadId  string          `json:"threadId"`
		Query     string          `json:"query"`
		Sequence  int             `json:"sequence"`
		CreatedAt string          `json:"createdAt"`
		UpdatedAt string          `json:"updatedAt"`
		Answers   []ThreadQAAResp `json:"answers"`
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

// request
type CustomerTIReq struct {
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

// request
type PATReq struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// request
type WorkspaceReq struct {
	Name string `json:"name"`
}

// request
type ThreadChatReq struct {
	Message string `json:"message"`
}

// response
type ThCustomerResp struct {
	CustomerId string
	Name       sql.NullString
}

func (c ThCustomerResp) MarshalJSON() ([]byte, error) {
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
type ThMemberResp struct {
	MemberId string
	Name     sql.NullString
}

func (m ThMemberResp) MarshalJSON() ([]byte, error) {
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

// response
type ThreadChatMessageResp struct {
	ThreadChatId        string
	ThreadChatMessageId string
	Body                string
	Sequence            int
	Customer            *ThCustomerResp
	Member              *ThMemberResp
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (thcresp ThreadChatMessageResp) MarshalJSON() ([]byte, error) {
	var customer *ThCustomerResp
	var member *ThMemberResp

	if thcresp.Customer != nil {
		customer = thcresp.Customer
	}

	if thcresp.Member != nil {
		member = thcresp.Member
	}

	aux := &struct {
		ThreadChatId        string          `json:"threadChatId"`
		ThreadChatMessageId string          `json:"threadChatMessageId"`
		Body                string          `json:"body"`
		Sequence            int             `json:"sequence"`
		Customer            *ThCustomerResp `json:"customer,omitempty"`
		Member              *ThMemberResp   `json:"member,omitempty"`
		CreatedAt           string          `json:"createdAt"`
		UpdatedAt           string          `json:"updatedAt"`
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
type ThreadChatResp struct {
	ThreadId  string
	Sequence  int
	Status    string
	Customer  ThCustomerResp
	Assignee  *ThMemberResp
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []ThreadChatMessageResp
}

func (thresp ThreadChatResp) MarshalJSON() ([]byte, error) {
	var assignee *ThMemberResp

	if thresp.Assignee != nil {
		assignee = thresp.Assignee
	}

	aux := &struct {
		ThreadId  string                  `json:"threadId"`
		Sequence  int                     `json:"sequence"`
		Status    string                  `json:"status"`
		Customer  ThCustomerResp          `json:"customer"`
		Assignee  *ThMemberResp           `json:"assignee"`
		CreatedAt string                  `json:"createdAt"`
		UpdatedAt string                  `json:"updatedAt"`
		Messages  []ThreadChatMessageResp `json:"messages"`
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

func AuthenticatePAT(ctx context.Context, db *pgxpool.Pool, token string) (Account, error) {
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

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

// AuthJWTClaims taken from Supabase JWT encoding
type AuthJWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func AuthenticateAccount(ctx context.Context, db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) (Account, error) {
	var account Account

	ath := r.Header.Get("Authorization")
	if ath == "" {
		return account, fmt.Errorf("cannot authenticate without authorization header")
	}

	cred := strings.Split(ath, " ")
	scheme := strings.ToLower(cred[0])

	if scheme == "token" {
		fmt.Println("authenticate via PAT")
		account, err := AuthenticatePAT(ctx, db, cred[1])
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
		ac, err := parseJWTToken(cred[1], []byte(hmacSecret))
		if err != nil {
			return account, fmt.Errorf("failed to parse JWT token with error: %v", err)
		}
		sub, err := ac.RegisteredClaims.GetSubject()
		if err != nil {
			return account, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}
		fmt.Printf("authenticated account with id: %s\n", sub)

		account, err := Account{}.GetByAuthUserId(ctx, db, sub)
		if err != nil {
			return account, fmt.Errorf("failed to get account by authUserId: %s with error: %v", ta.AuthUserId, err)
		}
		return account, nil
	} else {
		return account, fmt.Errorf("unsupported scheme: `%s` cannot authenticate", scheme)
	}
}

// Custom JWT claims for Customer
type CustomerJWTClaims struct {
	WorkspaceId string `json:"workspaceId"`
	ExternalId  string `json:"externalId"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	jwt.RegisteredClaims
}

func AuthenticateCustomer(ctx context.Context, db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) (Customer, error) {
	var customer Customer

	ath := r.Header.Get("Authorization")
	if ath == "" {
		return customer, fmt.Errorf("cannot authenticate without authorization header")
	}

	cred := strings.Split(ath, " ")
	scheme := strings.ToLower(cred[0])

	if scheme == "bearer" {
		fmt.Println("authenticate via JWTs")
		hmacSecret, err := GetEnv("SUPABASE_JWT_SECRET")
		if err != nil {
			return customer, fmt.Errorf("failed to get env SUPABASE_JWT_SECRET with error: %v", err)
		}
		cc, err := parseCustomerJWTToken(cred[1], []byte(hmacSecret))
		if err != nil {
			return customer, fmt.Errorf("failed to parse JWT token with error: %v", err)
		}
		sub, err := cc.RegisteredClaims.GetSubject()
		if err != nil {
			return customer, fmt.Errorf("cannot get subject from parsed token: %v", err)
		}
		fmt.Printf("authenticated customer with id: %s\n", sub)
		tc := Customer{WorkspaceId: cc.WorkspaceId, CustomerId: sub}
		customer, err := tc.GetById(ctx, db)
		if err != nil {
			return customer, fmt.Errorf("failed to get customer by customerId: %s with error: %v", tc.CustomerId, err)
		}
		return customer, nil
	} else {
		return customer, fmt.Errorf("unsupported scheme: `%s` cannot authenticate", scheme)
	}
}

func parseJWTToken(token string, hmacSecret []byte) (ac AuthJWTClaims, err error) {
	t, err := jwt.ParseWithClaims(token, &AuthJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return ac, fmt.Errorf("error validating jwt token with error: %v", err)
	} else if claims, ok := t.Claims.(*AuthJWTClaims); ok {
		return *claims, nil
		// sub, err := claims.RegisteredClaims.GetSubject()
		// if err != nil {
		// 	return auid, fmt.Errorf("cannot get subject from parsed token: %v", err)
		// }
	}

	return ac, fmt.Errorf("error parsing token: %v", token)
}

func parseCustomerJWTToken(token string, hmacSecret []byte) (cc CustomerJWTClaims, err error) {
	t, err := jwt.ParseWithClaims(token, &CustomerJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return cc, fmt.Errorf("error validating jwt token with error: %v", err)
	} else if claims, ok := t.Claims.(*CustomerJWTClaims); ok {
		return *claims, nil
	}
	return cc, fmt.Errorf("error parsing token: %v", token)
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
			ac, err := parseJWTToken(cred[1], []byte(hmacSecret))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			sub, err := ac.RegisteredClaims.GetSubject()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			name := ""
			account, isCreated, err := Account{}.GetOrCreateByAuthUserId(ctx, db, sub, ac.Email, DefaultAuthProvider, name)
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

func handleCreateAccountPAT(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var rb PATReq
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

		ap := AccountPAT{
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

func handleGetAccountPAT(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, err := AuthenticateAccount(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ap := AccountPAT{AccountId: account.AccountId}
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

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Header().Set("x-datetime", tm)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleCreateWorkspace(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var rb WorkspaceReq
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
		workspace := Workspace{AccountId: account.AccountId, Name: rb.Name}
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

// func handleLLMQueryEval(ctx context.Context, db *pgxpool.Pool) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer func(r io.ReadCloser) {
// 			_, _ = io.Copy(io.Discard, r)
// 			_ = r.Close()
// 		}(r.Body)

// 		var workspace Workspace

// 		var eval LLMRREval
// 		err := json.NewDecoder(r.Body).Decode(&eval)
// 		if err != nil {
// 			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 			return
// 		}

// 		workspaceId := r.PathValue("workspaceId")

// 		row, err := db.Query(ctx, `SELECT workspace_id, account_id,
// 			name, created_at, updated_at
// 			FROM workspace WHERE workspace_id = $1`,
// 			workspaceId)
// 		if err != nil {
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 			return
// 		}
// 		defer row.Close()

// 		if !row.Next() {
// 			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 			return
// 		}

// 		err = row.Scan(
// 			&workspace.WorkspaceId, &workspace.AccountId,
// 			&workspace.Name, &workspace.CreatedAt, &workspace.UpdatedAt)
// 		if err != nil {
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 			return
// 		}

// 		requestId := r.PathValue("requestId")

// 		_, err = db.Exec(ctx, `UPDATE llm_rr_log SET eval = $1
// 			WHERE workspace_id = $2 AND request_id = $3`,
// 			eval.Eval, workspace.WorkspaceId, requestId)
// 		if err != nil {
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusNoContent)
// 	})
// }

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

// TODO: work on this later - lot changes as we have update Thread Chat data model
func handleInitCustomerThreadQA(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var query ThreadQAReq

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

		workspace, err := Workspace{WorkspaceId: customer.WorkspaceId}.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		tq := ThreadQA{
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
		tqa := ThreadQAA{
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

		ans := make([]ThreadQAAResp, 0, 1)
		ans = append(ans, ThreadQAAResp{
			AnswerId: tqa.AnswerId,
			Answer:   tqa.Answer,
			Eval:     tqa.Eval,
			Sequence: tqa.Sequence,
		})
		resp := ThreadQAResp{
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

func handleInitCustomerThreadChat(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var message ThreadChatReq

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

		workspace, err := Workspace{WorkspaceId: customer.WorkspaceId}.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		th, thm, err := ThreadChat{}.CreateCustomerThChat(ctx, db, workspace, customer, message.Message)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		messages := make([]ThreadChatMessageResp, 0, 1)

		var msgCustomerRepr *ThCustomerResp
		var msgMemberRepr *ThMemberResp

		// for thread message - either of them
		if thm.CustomerId.Valid {
			msgCustomerRepr = &ThCustomerResp{
				CustomerId: thm.CustomerId.String,
				Name:       thm.CustomerName,
			}
		} else if thm.MemberId.Valid {
			msgMemberRepr = &ThMemberResp{
				MemberId: thm.MemberId.String,
				Name:     thm.MemberName,
			}
		}

		threadMessage := ThreadChatMessageResp{
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

		var threadAssigneeRepr *ThMemberResp

		// for thread
		threadCustomerRepr := ThCustomerResp{
			CustomerId: th.CustomerId,
			Name:       th.CustomerName,
		}

		// for thread
		if th.AssigneeId.Valid {
			threadAssigneeRepr = &ThMemberResp{
				MemberId: th.AssigneeId.String,
				Name:     th.AssigneeName,
			}
		}

		resp := ThreadChatResp{
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

		workspace, err := Workspace{WorkspaceId: customer.WorkspaceId}.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		ths, err := ThreadChat{}.GetListByWorkspaceCustomerId(ctx, db, workspace.WorkspaceId, customer.CustomerId)
		if err != nil {
			fmt.Printf("error in Get List By WorksapceId %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		threads := make([]ThreadChatResp, 0, 100)
		for _, th := range ths {
			messages := make([]ThreadChatMessageResp, 0, 1)

			var threadAssigneeRepr *ThMemberResp

			var msgCustomerRepr *ThCustomerResp
			var msgMemberRepr *ThMemberResp

			// for thread
			threadCustomerRepr := ThCustomerResp{
				CustomerId: th.ThreadChat.CustomerId,
				Name:       th.ThreadChat.CustomerName,
			}

			// for thread
			if th.ThreadChat.AssigneeId.Valid {
				threadAssigneeRepr = &ThMemberResp{
					MemberId: th.ThreadChat.AssigneeId.String,
					Name:     th.ThreadChat.AssigneeName,
				}
			}

			// for thread message - either of them
			if th.Message.CustomerId.Valid {
				msgCustomerRepr = &ThCustomerResp{
					CustomerId: th.Message.CustomerId.String,
					Name:       th.Message.CustomerName,
				}
			} else if th.Message.MemberId.Valid {
				msgMemberRepr = &ThMemberResp{
					MemberId: th.Message.MemberId.String,
					Name:     th.Message.MemberName,
				}
			}

			message := ThreadChatMessageResp{
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
			threads = append(threads, ThreadChatResp{
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

func handleCreateCustomerThreadChatMessage(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		threadId := r.PathValue("threadId")

		var message ThreadChatReq

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
		_, err = Workspace{WorkspaceId: customer.WorkspaceId}.GetById(ctx, db)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		th, err := ThreadChat{}.GetById(ctx, db, threadId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		thm, err := ThreadChatMessage{}.CreateCustomerThChatMessage(ctx, db, th, customer, message.Message)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var threadAssigneeRepr *ThMemberResp

		var msgCustomerRepr *ThCustomerResp
		var msgMemberRepr *ThMemberResp

		// for thread
		threadCustomerRepr := ThCustomerResp{
			CustomerId: th.CustomerId,
			Name:       th.CustomerName,
		}

		// for thread
		if th.AssigneeId.Valid {
			threadAssigneeRepr = &ThMemberResp{
				MemberId: th.AssigneeId.String,
				Name:     th.AssigneeName,
			}
		}

		// for thread message - either of them
		if thm.CustomerId.Valid {
			msgCustomerRepr = &ThCustomerResp{
				CustomerId: thm.CustomerId.String,
				Name:       thm.CustomerName,
			}
		} else if thm.MemberId.Valid {
			msgMemberRepr = &ThMemberResp{
				MemberId: thm.MemberId.String,
				Name:     thm.MemberName,
			}
		}

		threadMessage := ThreadChatMessageResp{
			ThreadChatId:        th.ThreadChatId,
			ThreadChatMessageId: thm.ThreadChatMessageId,
			Body:                thm.Body,
			Sequence:            thm.Sequence,
			Customer:            msgCustomerRepr,
			Member:              msgMemberRepr,
			CreatedAt:           thm.CreatedAt,
			UpdatedAt:           thm.UpdatedAt,
		}

		messages := make([]ThreadChatMessageResp, 0, 1)
		messages = append(messages, threadMessage)
		resp := ThreadChatResp{
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

func handleCreateMemberThreadChatMessage(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		workspaceId := r.PathValue("workspaceId")
		threadId := r.PathValue("threadId")

		var message ThreadChatReq

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

		member, err := Member{}.GetWorkspaceMemberByAccountId(ctx, db, workspaceId, account.AccountId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		th, err := ThreadChat{}.GetById(ctx, db, threadId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		thm, err := ThreadChatMessage{}.CreateMemberThChatMessage(ctx, db, th, member, message.Message)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !th.AssigneeId.Valid {
			fmt.Println("Thread Chat not yet assigned will try to assign Member...")
			thAssigned, err := ThreadChat{}.AssignMember(ctx, db, member.MemberId, th.ThreadChatId)
			if err != nil {
				fmt.Printf("(silent) failed to assign member to Thread Chat with error: %v\n", err)
			} else {
				th = thAssigned
			}
		}

		var threadAssigneeRepr *ThMemberResp

		var msgCustomerRepr *ThCustomerResp
		var msgMemberRepr *ThMemberResp

		// for thread
		threadCustomerRepr := ThCustomerResp{
			CustomerId: th.CustomerId,
			Name:       th.CustomerName,
		}

		// for thread
		if th.AssigneeId.Valid {
			threadAssigneeRepr = &ThMemberResp{
				MemberId: th.AssigneeId.String,
				Name:     th.AssigneeName,
			}
		}

		// for thread message - either of them
		if thm.CustomerId.Valid {
			msgCustomerRepr = &ThCustomerResp{
				CustomerId: thm.CustomerId.String,
				Name:       thm.CustomerName,
			}
		} else if thm.MemberId.Valid {
			msgMemberRepr = &ThMemberResp{
				MemberId: thm.MemberId.String,
				Name:     thm.MemberName,
			}
		}

		threadMessage := ThreadChatMessageResp{
			ThreadChatId:        th.ThreadChatId,
			ThreadChatMessageId: thm.ThreadChatMessageId,
			Body:                thm.Body,
			Sequence:            thm.Sequence,
			Customer:            msgCustomerRepr,
			Member:              msgMemberRepr,
			CreatedAt:           thm.CreatedAt,
			UpdatedAt:           thm.UpdatedAt,
		}

		messages := make([]ThreadChatMessageResp, 0, 1)
		messages = append(messages, threadMessage)
		resp := ThreadChatResp{
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

func handleGetCustomerThreadChatMessages(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := AuthenticateCustomer(ctx, db, w, r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		threadId := r.PathValue("threadId")

		th, err := ThreadChat{}.GetById(ctx, db, threadId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		results, err := ThreadChatMessage{}.GetListByThreadId(ctx, db, th.ThreadChatId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		messages := make([]ThreadChatMessageResp, 0, 100)
		for _, thm := range results {
			var msgCustomerRepr *ThCustomerResp
			var msgMemberRepr *ThMemberResp

			// for thread message - either of them
			if thm.CustomerId.Valid {
				msgCustomerRepr = &ThCustomerResp{
					CustomerId: thm.CustomerId.String,
					Name:       thm.CustomerName,
				}
			} else if thm.MemberId.Valid {
				msgMemberRepr = &ThMemberResp{
					MemberId: thm.MemberId.String,
					Name:     thm.MemberName,
				}
			}

			threadMessage := ThreadChatMessageResp{
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

		var threadAssigneeRepr *ThMemberResp

		// for thread
		threadCustomerRepr := ThCustomerResp{
			CustomerId: th.CustomerId,
			Name:       th.CustomerName,
		}

		// for thread
		if th.AssigneeId.Valid {
			threadAssigneeRepr = &ThMemberResp{
				MemberId: th.AssigneeId.String,
				Name:     th.AssigneeName,
			}
		}

		resp := ThreadChatResp{
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

func handleCustomerTokenIssue(ctx context.Context, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		var rb CustomerTIReq
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

		tw := Workspace{WorkspaceId: workspaceId, AccountId: account.AccountId}
		workspace, err := tw.GetAccountWorkspace(ctx, db)
		if err != nil {
			fmt.Printf("failed to get account workspace or does not exist with error: %v\n", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		tc := Customer{
			WorkspaceId: workspace.WorkspaceId,
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
				jwt, err := customer.MakeJWT()
				if err != nil {
					fmt.Printf("failed to make jwt token with error: %v\n", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				resp := CustomerTIResp{
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
				jwt, err := customer.MakeJWT()
				if err != nil {
					fmt.Printf("failed to make jwt token with error: %v\n", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				resp := CustomerTIResp{
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
				jwt, err := customer.MakeJWT()
				if err != nil {
					fmt.Printf("failed to make jwt token with error: %v\n", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				resp := CustomerTIResp{
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
			var customer Customer
			fmt.Printf("based on identifiers check for customer in workspaceId: %v\n", workspaceId)
			if tc.ExternalId.Valid {
				extId := tc.ExternalId.String
				fmt.Printf("get customer by externalId %s\n", extId)
				customer, err = customer.GetByExtId(ctx, db, workspaceId, extId)
				if err != nil {
					fmt.Printf("failed to get customer by externalId %s with error: %v\n", extId, err)
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
				resp := CustomerTIResp{
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
			} else if tc.Email.Valid {
				email := tc.Email.String
				fmt.Printf("get customer by email %s\n", email)
				customer, err = customer.GetByEmail(ctx, db, workspaceId, email)
				if err != nil {
					fmt.Printf("failed to get customer by email %s with error: %v\n", email, err)
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
				resp := CustomerTIResp{
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
			} else if tc.Phone.Valid {
				phone := tc.Phone.String
				fmt.Printf("get customer by phone %s\n", phone)
				customer, err = customer.GetByPhone(ctx, db, workspaceId, phone)
				if err != nil {
					fmt.Printf("failed to get customer by phone %s with error: %v\n", phone, err)
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
				resp := CustomerTIResp{
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
	mux.HandleFunc("GET /{$}", handleGetIndex)

	// web
	mux.Handle("POST /accounts/auth/{$}",
		handleAuthAccountMaker(ctx, db))

	// sdk+web
	mux.Handle("POST /pats/{$}",
		handleCreateAccountPAT(ctx, db))

	// sdk+web
	mux.Handle("GET /pats/{$}",
		handleGetAccountPAT(ctx, db))

	// sdk+web
	mux.Handle("POST /workspaces/{$}",
		handleCreateWorkspace(ctx, db))

	// sdk+web
	mux.Handle("GET /workspaces/{$}",
		handleGetWorkspaces(ctx, db))

	// customer
	mux.Handle("GET /-/me/{$}", handleGetCustomer(ctx, db))

	// customer
	mux.Handle("POST /-/threads/qa/{$}",
		handleInitCustomerThreadQA(ctx, db))

	// customer
	mux.Handle("POST /-/threads/chat/{$}",
		handleInitCustomerThreadChat(ctx, db))

	// customer
	mux.Handle("POST /-/threads/chat/{threadId}/messages/{$}",
		handleCreateCustomerThreadChatMessage(ctx, db))

	// customer
	mux.Handle("GET /-/threads/chat/{$}",
		handleGetCustomerThreadChats(ctx, db))

	// customer
	mux.Handle("GET /-/threads/chat/{threadId}/messages/{$}",
		handleGetCustomerThreadChatMessages(ctx, db))

	// sdk+web
	mux.Handle("POST /workspaces/{workspaceId}/threads/chat/{threadId}/messages/{$}",
		handleCreateMemberThreadChatMessage(ctx, db))

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
