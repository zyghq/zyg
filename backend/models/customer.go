package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/rs/xid"
	"github.com/sanchitrk/namingo"
	"github.com/zyghq/zyg"
	"strings"
	"time"
)

type Customer struct {
	WorkspaceId     string
	CustomerId      string
	ExternalId      sql.NullString
	Email           sql.NullString
	Phone           sql.NullString
	Name            string
	IsEmailVerified bool
	Role            string
	UpdatedAt       time.Time
	CreatedAt       time.Time
}

func (c Customer) GenId() string {
	return "cs" + xid.New().String()
}

func (c Customer) Visitor() string {
	return "visitor"
}

func (c Customer) Lead() string {
	return "lead"
}

func (c Customer) Engaged() string {
	return "engaged"
}

func (c Customer) IsVisitor() bool {
	return c.Role == c.Visitor()
}

func (c Customer) AnonName() string {
	return namingo.Generate(2, " ", namingo.TitleCase())
}

func (c Customer) AvatarUrl() string {
	url := zyg.GetAvatarBaseURL()
	// url may or may not have a trailing slash
	// add a trailing slash if it doesn't have one
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return url + c.CustomerId
}

func (c Customer) AsCustomerActor() CustomerActor {
	return CustomerActor{
		CustomerId: c.CustomerId,
		Name:       c.Name,
	}
}

// IdentityHash is a hash of the customer's identity
// Combined these fields create a unique hash for the customer
// (XXX): You might have to update this if you plan to add more identity fields
func (c Customer) IdentityHash() string {
	h := sha256.New()
	// Combine all fields into a single string
	identityString := fmt.Sprintf("%s:%s:%s:%s:%s:%t",
		c.WorkspaceId,
		c.CustomerId,
		c.ExternalId.String,
		c.Email.String,
		c.Phone.String,
		c.IsEmailVerified,
	)

	// Write the combined string to the hash
	h.Write([]byte(identityString))

	// Return the hash as a base64 encoded string
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (c Customer) MakeCopy() Customer {
	return Customer{
		WorkspaceId:     c.WorkspaceId,
		CustomerId:      c.CustomerId,
		ExternalId:      c.ExternalId,
		Email:           c.Email,
		Phone:           c.Phone,
		Name:            c.Name,
		IsEmailVerified: c.IsEmailVerified,
		Role:            c.Role,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}

func (c Customer) HasEmail() bool {
	return c.Email.Valid
}
