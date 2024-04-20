package model

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"

	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/internal/auth"
)

type mErr string

func (err mErr) Error() string {
	return string(err)
}

const (
	ErrNothing   = mErr("nothing found")
	ErrMapping   = mErr("failed to map")
	ErrQuery     = mErr("failed to query")
	ErrSomething = mErr("something went wrong")
)

func GenToken(length int, prefix string) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return prefix + base64.URLEncoding.EncodeToString(buffer)[:length], nil
}

type ThreadStatus struct{}

func (s ThreadStatus) Todo() string {
	return "todo"
}

func (s ThreadStatus) Done() string {
	return "done"
}

func (s ThreadStatus) InProgress() string {
	return "inprogress"
}

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

func (a Account) GetOrCreateByAuthUserId(ctx context.Context, db *pgxpool.Pool) (Account, bool, error) {
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

	row, err := db.Query(ctx, st, aId, a.AuthUserId, a.Email, a.Provider, a.Name)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return a, isCreated, ErrQuery
	}
	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return a, isCreated, ErrNothing
	}

	err = row.Scan(
		&a.AccountId, &a.AuthUserId,
		&a.Email, &a.Provider, &a.Name,
		&a.CreatedAt, &a.UpdatedAt,
		&isCreated,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return a, isCreated, ErrMapping
	}

	return a, isCreated, nil
}

func (a Account) GetByAuthUserId(ctx context.Context, db *pgxpool.Pool) (Account, error) {
	row, err := db.Query(ctx, `SELECT 
		account_id, auth_user_id, email,
		provider, name, created_at, updated_at
		FROM account WHERE auth_user_id = $1`, a.AuthUserId)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return a, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return a, ErrNothing
	}

	err = row.Scan(
		&a.AccountId, &a.AuthUserId,
		&a.Email, &a.Provider,
		&a.Name, &a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return a, ErrMapping
	}

	return a, nil
}

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
		fmt.Printf("failed to generate token got error: %v\n", err)
		return ap, ErrSomething
	}
	stmt := `INSERT INTO account_pat(account_id, pat_id, token, name, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING account_id, pat_id, token, name, description, created_at, updated_at`

	row, err := db.Query(ctx, stmt, ap.AccountId, apId, token, ap.Name, ap.Description)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return ap, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return ap, ErrNothing
	}

	err = row.Scan(
		&ap.AccountId, &ap.PatId, &ap.Token,
		&ap.Name, &ap.Description, &ap.CreatedAt, &ap.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return ap, ErrMapping
	}
	return ap, nil
}

func (ap AccountPAT) GetListByAccountId(ctx context.Context, db *pgxpool.Pool) ([]AccountPAT, error) {
	aps := make([]AccountPAT, 0, 100)

	stmt := `SELECT account_id, pat_id, token, name, description,
		created_at, updated_at
		FROM account_pat WHERE account_id = $1
		ORDER BY created_at DESC LIMIT 100`

	rows, err := db.Query(ctx, stmt, ap.AccountId)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return aps, ErrQuery
	}

	defer rows.Close()

	if !rows.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return aps, ErrNothing
	}

	for rows.Next() {
		var ap AccountPAT
		err = rows.Scan(
			&ap.AccountId, &ap.PatId, &ap.Token,
			&ap.Name, &ap.Description, &ap.CreatedAt, &ap.UpdatedAt,
		)
		if err != nil {
			fmt.Printf("failed to scan got error: %v", err)
			return aps, ErrMapping
		}
		aps = append(aps, ap)
	}

	return aps, nil
}

func (ap AccountPAT) GetByToken(ctx context.Context, db *pgxpool.Pool) (Account, error) {
	var account Account

	stmt := `SELECT
		a.account_id, a.email,
		a.provider, a.auth_user_id, a.name,
		a.created_at, a.updated_at
		FROM account a
		INNER JOIN account_pat ap ON a.account_id = ap.account_id
		WHERE ap.token = $1`

	row, err := db.Query(ctx, stmt, ap.Token)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return account, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return account, ErrNothing
	}

	err = row.Scan(
		&account.AccountId, &account.Email,
		&account.Provider, &account.AuthUserId, &account.Name,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return account, ErrMapping
	}

	return account, nil
}

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
	var member Member

	tx, err := db.Begin(ctx)
	if err != nil {
		fmt.Printf("failed to start db transaction got error: %v\n", err)
		return w, ErrQuery
	}
	defer tx.Rollback(ctx)

	wId := w.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO workspace(workspace_id, account_id, name)
		VALUES ($1, $2, $3)
		RETURNING
		workspace_id, account_id, name, created_at, updated_at`, wId, w.AccountId, w.Name).Scan(
		&w.WorkspaceId, &w.AccountId,
		&w.Name, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to insert got error: %v\n", err)
		return w, ErrQuery
	}

	mId := member.GenId()
	err = tx.QueryRow(ctx, `INSERT INTO member(workspace_id, account_id, member_id, role)
		VALUES ($1, $2, $3, $4)
		RETURNING
		workspace_id, account_id, member_id, name, role, created_at, updated_at`,
		wId, w.AccountId, mId, "primary").Scan(
		&member.WorkspaceId, &member.AccountId,
		&member.MemberId, &member.Name, &member.Role,
		&member.CreatedAt, &member.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to insert got error: %v\n", err)
		return w, ErrQuery
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit db transaction got error: %v\n", err)
		return w, ErrQuery
	}
	return w, nil
}

func (w Workspace) GetAccountWorkspace(ctx context.Context, db *pgxpool.Pool) (Workspace, error) {
	row, err := db.Query(ctx, `SELECT 
		workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE account_id = $1 AND workspace_id = $2`, w.AccountId, w.WorkspaceId)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return w, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return w, ErrNothing
	}

	err = row.Scan(
		&w.WorkspaceId, &w.AccountId,
		&w.Name, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return w, ErrMapping
	}
	return w, nil
}

func (w Workspace) GetById(ctx context.Context, db *pgxpool.Pool) (Workspace, error) {
	row, err := db.Query(ctx, `SELECT 
		workspace_id, account_id, name, created_at, updated_at
		FROM workspace WHERE workspace_id = $1`, w.WorkspaceId)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return w, ErrQuery
	}
	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return w, ErrNothing
	}

	err = row.Scan(
		&w.WorkspaceId, &w.AccountId,
		&w.Name, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return w, ErrMapping
	}
	return w, nil
}

type Member struct {
	WorkspaceId string
	AccountId   string
	MemberId    string
	Name        sql.NullString
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

func (m Member) GenId() string {
	return "m_" + xid.New().String()
}

func (m Member) GetWorkspaceMemberByAccountId(ctx context.Context, db *pgxpool.Pool) (Member, error) {
	row, err := db.Query(ctx, `SELECT
		workspace_id, account_id, member_id, name, role, created_at, updated_at
		FROM member WHERE workspace_id = $1 AND account_id = $2`, m.WorkspaceId, m.AccountId)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return m, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return m, ErrNothing
	}

	err = row.Scan(
		&m.WorkspaceId, &m.AccountId,
		&m.MemberId, &m.Name, &m.Role,
		&m.CreatedAt, &m.UpdatedAt,
	)

	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return m, ErrMapping
	}

	return m, nil
}

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

func (c Customer) GetWrkCustomerById(ctx context.Context, db *pgxpool.Pool) (Customer, error) {
	row, err := db.Query(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND customer_id = $2`, c.WorkspaceId, c.CustomerId)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return c, ErrQuery
	}
	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return c, ErrNothing
	}

	err = row.Scan(
		&c.WorkspaceId, &c.CustomerId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.Name,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return c, ErrMapping
	}
	return c, nil
}

func (c Customer) GetWrkCustomerByExtId(ctx context.Context, db *pgxpool.Pool) (Customer, error) {
	row, err := db.Query(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND external_id = $2`, c.WorkspaceId, c.ExternalId)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return c, ErrQuery
	}
	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return c, ErrNothing
	}

	err = row.Scan(
		&c.WorkspaceId, &c.CustomerId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.Name,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return c, ErrMapping
	}
	return c, nil
}

func (c Customer) GetWrkCustomerByEmail(ctx context.Context, db *pgxpool.Pool) (Customer, error) {
	row, err := db.Query(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND email = $2`, c.WorkspaceId, c.Email)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return c, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return c, ErrNothing
	}

	err = row.Scan(
		&c.WorkspaceId, &c.CustomerId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.Name,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return c, ErrMapping
	}

	return c, nil
}

func (c Customer) GetWrkCustomerByPhone(ctx context.Context, db *pgxpool.Pool) (Customer, error) {
	row, err := db.Query(ctx, `SELECT 
		workspace_id, customer_id,
		external_id, email,
		phone, name, created_at, updated_at
		FROM customer WHERE workspace_id = $1 AND phone = $2`, c.WorkspaceId, c.Phone)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return c, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return c, ErrNothing
	}

	err = row.Scan(
		&c.WorkspaceId, &c.CustomerId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.Name,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return c, ErrMapping
	}
	return c, nil
}

func (c Customer) GetOrCreateWrkCustomerByExtId(ctx context.Context, db *pgxpool.Pool) (Customer, bool, error) {
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

	var isCreated bool
	row, err := db.Query(ctx, st, cId, c.WorkspaceId, c.ExternalId, c.Email, c.Phone)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return c, isCreated, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return c, isCreated, ErrNothing
	}

	err = row.Scan(
		&c.CustomerId, &c.WorkspaceId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.CreatedAt,
		&c.UpdatedAt, &isCreated,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return c, isCreated, ErrMapping
	}
	return c, isCreated, nil
}

func (c Customer) GetOrCreateWrkCustomerByEmail(ctx context.Context, db *pgxpool.Pool) (Customer, bool, error) {
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

	var isCreated bool
	row, err := db.Query(ctx, st, cId, c.WorkspaceId, c.ExternalId, c.Email, c.Phone)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return c, isCreated, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return c, isCreated, ErrNothing
	}

	err = row.Scan(
		&c.CustomerId, &c.WorkspaceId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.CreatedAt,
		&c.UpdatedAt, &isCreated,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return c, isCreated, ErrMapping
	}
	return c, isCreated, nil
}

func (c Customer) GetOrCreateWrkCustomerByPhone(ctx context.Context, db *pgxpool.Pool) (Customer, bool, error) {
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

	var isCreated bool
	row, err := db.Query(ctx, st, cId, c.WorkspaceId, c.ExternalId, c.Email, c.Phone)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return c, isCreated, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return c, isCreated, ErrNothing
	}

	err = row.Scan(
		&c.CustomerId, &c.WorkspaceId,
		&c.ExternalId, &c.Email,
		&c.Phone, &c.CreatedAt,
		&c.UpdatedAt, &isCreated,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return c, isCreated, ErrMapping
	}
	return c, isCreated, nil
}

func (c Customer) MakeJWT() (string, error) {

	var externalId string
	var email string
	var phone string

	audience := []string{"customer"}

	sk, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
	if err != nil {
		return "", fmt.Errorf("failed to get env SUPABASE_JWT_SECRET got error: %v", err)
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

	claims := auth.CustomerJWTClaims{
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
		return "", fmt.Errorf("failed to sign JWT token got error: %v", err)
	}
	return jwt, nil
}

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
		fmt.Printf("failed to start db transaction got error: %v\n", err)
		return th, thm, ErrQuery
	}

	defer tx.Rollback(ctx)

	thId := th.GenId()
	th.Status = ThreadStatus{}.Todo()
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
		fmt.Printf("failed to insert in thread chat got error: %v\n", err)
		return th, thm, ErrQuery
	}

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
		fmt.Printf("failed to insert thread chat message got error: %v\n", err)
		return th, thm, ErrQuery
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit db transaction got error: %v\n", err)
		return th, thm, ErrQuery
	}

	// set customer attributes we already have
	th.CustomerName = c.Name
	thm.CustomerName = c.Name
	return th, thm, nil
}

func (th ThreadChat) GetById(ctx context.Context, db *pgxpool.Pool) (ThreadChat, error) {
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

	err := db.QueryRow(ctx, stmt, th.ThreadChatId).Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return th, ErrQuery
	}

	return th, nil
}

func (th ThreadChat) GetListByWorkspaceCustomerId(ctx context.Context, db *pgxpool.Pool) ([]ThreadChatWithMessage, error) {
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

	rows, err := db.Query(ctx, stmt, th.WorkspaceId, th.CustomerId)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return ths, ErrQuery
	}

	defer rows.Close()

	if !rows.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return ths, ErrNothing
	}

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
			fmt.Printf("failed to scan got error: %v", err)
			return ths, ErrMapping
		}

		thm := ThreadChatWithMessage{
			ThreadChat: th,
			Message:    tc,
		}
		ths = append(ths, thm)
	}

	return ths, nil
}

func (th ThreadChat) AssignMember(ctx context.Context, db *pgxpool.Pool) (ThreadChat, error) {
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

	row, err := db.Query(ctx, stmt, th.AssigneeId, th.ThreadChatId)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return th, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return th, ErrNothing
	}

	err = row.Scan(
		&th.WorkspaceId, &th.CustomerId, &th.CustomerName,
		&th.AssigneeId, &th.AssigneeName,
		&th.ThreadChatId, &th.Title, &th.Summary,
		&th.Sequence, &th.Status, &th.CreatedAt, &th.UpdatedAt,
	)

	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return th, ErrMapping
	}

	return th, nil
}

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
	ctx context.Context, db *pgxpool.Pool, c Customer, msg string,
) (ThreadChatMessage, error) {
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

	err := db.QueryRow(ctx, stmt, thm.ThreadChatId, thmId, msg, c.CustomerId).Scan(
		&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
		&thm.Sequence, &thm.CustomerId, &thm.CustomerName,
		&thm.MemberId, &thm.MemberName,
		&thm.CreatedAt, &thm.UpdatedAt,
	)

	if err != nil {
		fmt.Printf("failed to insert got error: %v\n", err)
		return thm, ErrQuery
	}
	return thm, nil
}

func (thm ThreadChatMessage) CreateMemberThChatMessage(
	ctx context.Context, db *pgxpool.Pool, m Member, msg string,
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
	err := db.QueryRow(ctx, stmt, thm.ThreadChatId, thmId, msg, m.MemberId).Scan(
		&thm.ThreadChatId, &thm.ThreadChatMessageId, &thm.Body,
		&thm.Sequence, &thm.CustomerId, &thm.CustomerName,
		&thm.MemberId, &thm.MemberName,
		&thm.CreatedAt, &thm.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed insert for thread chat message for thread id %s with error: %v\n", thm.ThreadChatId, err)
		return thm, err
	}
	return thm, nil
}

func (thm ThreadChatMessage) GetListByThreadChatId(ctx context.Context, db *pgxpool.Pool) ([]ThreadChatMessage, error) {
	messages := make([]ThreadChatMessage, 0, 100)
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

	rows, err := db.Query(ctx, stmt, thm.ThreadChatId)
	if err != nil {
		fmt.Printf("failed to query got error: %v\n", err)
		return messages, ErrQuery
	}

	defer rows.Close()

	if !rows.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return messages, ErrNothing
	}

	for rows.Next() {
		var m ThreadChatMessage
		err = rows.Scan(
			&m.ThreadChatId, &m.ThreadChatMessageId, &m.Body, &m.Sequence,
			&m.CreatedAt, &m.UpdatedAt, &m.CustomerId, &m.CustomerName,
			&m.MemberId, &m.MemberName,
		)
		if err != nil {
			fmt.Printf("failed to scan got error: %v", err)
			return messages, ErrMapping
		}
		messages = append(messages, m)
	}

	return messages, nil
}

type ThreadChatWithMessage struct {
	ThreadChat ThreadChat
	Message    ThreadChatMessage
}

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
		fmt.Printf("failed to query got error: %v\n", err)
		return thread, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return thread, ErrNothing
	}

	err = row.Scan(
		&thread.WorkspaceId, &thread.ThreadQAId, &thread.AnswerId, &thread.Answer,
		&thread.Eval, &thread.Sequence, &thread.CreatedAt, &thread.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return thread, ErrMapping
	}

	return thread, nil
}

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
		fmt.Printf("failed to query got error: %v\n", err)
		return thread, ErrQuery
	}

	defer row.Close()

	if !row.Next() {
		fmt.Printf("no rows got error: %v\n", sql.ErrNoRows)
		return thread, ErrNothing
	}

	err = row.Scan(
		&thread.WorkspaceId, &thread.CustomerId, &thread.ThreadId, &thread.ParentThreadId,
		&thread.Query, &thread.Title, &thread.Summary, &thread.Sequence,
		&thread.CreatedAt, &thread.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("failed to scan got error: %v", err)
		return thread, ErrMapping
	}

	return thread, nil
}
