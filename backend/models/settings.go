package models

import (
	"encoding/json"
	"time"
)

type PostmarkMailServerSetting struct {
	WorkspaceId              string     `json:"workspaceId"`
	ServerId                 int64      `json:"serverId"`
	ServerToken              string     `json:"serverToken"`
	IsEnabled                bool       `json:"isEnabled"`
	Email                    string     `json:"email"`
	Domain                   string     `json:"domain"`
	HasError                 bool       `json:"hasError"`
	InboundEmail             *string    `json:"inboundEmail"` // After creating the server
	HasForwardingEnabled     bool       `json:"hasForwardingEnabled"`
	HasDNS                   bool       `json:"hasDNS"`
	IsDNSVerified            bool       `json:"isDNSVerified"`
	DNSVerifiedAt            *time.Time `json:"dnsVerifiedAt"`            // After DNS is verified
	DNSDomainId              *int64     `json:"dnsDomainId"`              // After adding domain
	DKIMPendingHost          *string    `json:"dkimPendingHost"`          // After adding domain
	DKIMPendingTextValue     *string    `json:"dkimPendingTextValue"`     // After adding domain
	DKIMUpdateStatus         *string    `json:"dkimUpdateStatus"`         // After adding domain
	ReturnPathDomain         *string    `json:"returnPathDomain"`         // After adding domain
	ReturnPathDomainCNAME    *string    `json:"returnPathDomainCNAME"`    // After adding domain
	ReturnPathDomainVerified bool       `json:"returnPathDomainVerified"` // After adding domain
	CreatedAt                time.Time  `json:"createdAt"`
	UpdatedAt                time.Time  `json:"updatedAt"`
}

func (pm PostmarkMailServerSetting) MarshalJSON() ([]byte, error) {
	type Aux PostmarkMailServerSetting
	var token string
	maskLeft := func(s string) string {
		rs := []rune(s)
		for i := range rs[:len(rs)-4] {
			rs[i] = '*'
		}
		return string(rs)
	}

	token = maskLeft(pm.ServerToken)
	aux := &struct {
		Aux
		ServerToken string `json:"serverToken"`
	}{
		Aux:         Aux(pm),
		ServerToken: token,
	}
	return json.Marshal(aux)
}
