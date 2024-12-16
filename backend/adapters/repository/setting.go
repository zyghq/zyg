package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/cristalhq/builq"
	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/models"
	"log/slog"
)

func postmarkMailServerSettingCols() builq.Columns {
	return builq.Columns{
		"workspace_id",
		"server_id",
		"server_token",
		"is_enabled",
		"email",
		"domain",
		"has_error",
		"inbound_email", // nullable
		"has_forwarding_enabled",
		"has_dns",
		"is_dns_verified",
		"dns_verified_at",          // nullable
		"dns_domain_id",            // nullable
		"dkim_host",                // nullable
		"dkim_text_value",          // nullable
		"dkim_update_status",       // nullable
		"return_path_domain",       // nullable
		"return_path_domain_cname", // nullable
		"return_path_domain_verified",
		"created_at",
		"updated_at",
	}
}

func (wrk *WorkspaceDB) SavePostmarkMailServerSetting(
	ctx context.Context, setting models.PostmarkMailServerSetting) (models.PostmarkMailServerSetting, error) {
	q := builq.New()
	cols := postmarkMailServerSettingCols()

	var (
		inboundEmail  sql.NullString
		dnsVerifiedAt sql.NullTime
		dnsDomainId   sql.NullInt64
		dkimHost, dkimTextValue, dkimUpdateStatus,
		returnPathDomain, returnPathDomainCname sql.NullString
	)

	if setting.InboundEmail != nil {
		inboundEmail = sql.NullString{
			String: *setting.InboundEmail,
			Valid:  true,
		}
	}
	if setting.DNSVerifiedAt != nil {
		dnsVerifiedAt = sql.NullTime{
			Time:  *setting.DNSVerifiedAt,
			Valid: true,
		}
	}
	if setting.DNSDomainId != nil {
		dnsDomainId = sql.NullInt64{
			Int64: *setting.DNSDomainId,
			Valid: true,
		}
	}

	if setting.DKIMHost != nil {
		dkimHost = sql.NullString{
			String: *setting.DKIMHost,
			Valid:  true,
		}
	}

	if setting.DKIMTextValue != nil {
		dkimTextValue = sql.NullString{
			String: *setting.DKIMTextValue,
			Valid:  true,
		}
	}

	if setting.DKIMUpdateStatus != nil {
		dkimUpdateStatus = sql.NullString{
			String: *setting.DKIMUpdateStatus,
			Valid:  true,
		}
	}

	if setting.ReturnPathDomain != nil {
		returnPathDomain = sql.NullString{
			String: *setting.ReturnPathDomain,
			Valid:  true,
		}
	}

	if setting.ReturnPathDomainCNAME != nil {
		returnPathDomainCname = sql.NullString{
			String: *setting.ReturnPathDomainCNAME,
			Valid:  true,
		}
	}

	insertParams := []any{
		setting.WorkspaceId, setting.ServerId, setting.ServerToken, setting.IsEnabled, setting.Email,
		setting.Domain,
		setting.HasError, inboundEmail,
		setting.HasForwardingEnabled,
		setting.HasDNS, setting.IsDNSVerified, dnsVerifiedAt, dnsDomainId,
		dkimHost, dkimTextValue, dkimUpdateStatus,
		returnPathDomain, returnPathDomainCname, setting.ReturnPathDomainVerified,
		setting.CreatedAt,
		setting.UpdatedAt,
	}

	q("INSERT INTO postmark_mail_server_setting (%s)", cols)
	q("VALUES (%+$)", insertParams)
	q("ON CONFLICT (workspace_id) DO UPDATE SET")
	q("server_id = EXCLUDED.server_id,")
	q("server_token = EXCLUDED.server_token,")
	q("is_enabled = EXCLUDED.is_enabled,")
	q("email = EXCLUDED.email,")
	q("domain = EXCLUDED.domain,")
	q("has_error = EXCLUDED.has_error,")
	q("inbound_email = EXCLUDED.inbound_email,")
	q("has_forwarding_enabled = EXCLUDED.has_forwarding_enabled,")
	q("has_dns = EXCLUDED.has_dns,")
	q("is_dns_verified = EXCLUDED.is_dns_verified,")
	q("dns_verified_at = EXCLUDED.dns_verified_at,")
	q("dns_domain_id = EXCLUDED.dns_domain_id,")
	q("dkim_host = EXCLUDED.dkim_host,")
	q("dkim_text_value = EXCLUDED.dkim_text_value,")
	q("dkim_update_status = EXCLUDED.dkim_update_status,")
	q("return_path_domain = EXCLUDED.return_path_domain,")
	q("return_path_domain_cname = EXCLUDED.return_path_domain_cname,")
	q("return_path_domain_verified = EXCLUDED.return_path_domain_verified,")
	q("updated_at = NOW()")
	q("RETURNING %s", cols)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.PostmarkMailServerSetting{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = wrk.db.QueryRow(ctx, stmt, insertParams...).Scan(
		&setting.WorkspaceId, &setting.ServerId, &setting.ServerToken, &setting.IsEnabled, &setting.Email,
		&setting.Domain,
		&setting.HasError, &inboundEmail,
		&setting.HasForwardingEnabled,
		&setting.HasDNS, &setting.IsDNSVerified, &dnsVerifiedAt, &dnsDomainId,
		&dkimHost, &dkimTextValue, &dkimUpdateStatus,
		&returnPathDomain, &returnPathDomainCname, &setting.ReturnPathDomainVerified,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.PostmarkMailServerSetting{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return models.PostmarkMailServerSetting{}, ErrQuery
	}

	if inboundEmail.Valid {
		setting.InboundEmail = &inboundEmail.String
	}
	if dnsVerifiedAt.Valid {
		setting.DNSVerifiedAt = &dnsVerifiedAt.Time
	}
	if dnsDomainId.Valid {
		setting.DNSDomainId = &dnsDomainId.Int64
	}
	if dkimHost.Valid {
		setting.DKIMHost = &dkimHost.String
	}
	if dkimTextValue.Valid {
		setting.DKIMTextValue = &dkimTextValue.String
	}
	if dkimUpdateStatus.Valid {
		setting.DKIMUpdateStatus = &dkimUpdateStatus.String
	}
	if returnPathDomain.Valid {
		setting.ReturnPathDomain = &returnPathDomain.String
	}
	if returnPathDomainCname.Valid {
		setting.ReturnPathDomainCNAME = &returnPathDomainCname.String
	}
	return setting, nil
}

func (wrk *WorkspaceDB) FetchPostmarkMailServerSettingByWorkspaceId(
	ctx context.Context, workspaceId string) (models.PostmarkMailServerSetting, error) {
	var setting models.PostmarkMailServerSetting
	var (
		inboundEmail  sql.NullString
		dnsVerifiedAt sql.NullTime
		dnsDomainId   sql.NullInt64
		dkimHost, dkimTextValue, dkimUpdateStatus,
		returnPathDomain, returnPathDomainCname sql.NullString
	)

	q := builq.New()
	cols := postmarkMailServerSettingCols()

	q("SELECT %s FROM postmark_mail_server_setting", cols)
	q("WHERE workspace_id = %$", workspaceId)

	stmt, _, err := q.Build()
	if err != nil {
		slog.Error("failed to build query", slog.Any("err", err))
		return models.PostmarkMailServerSetting{}, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := q.DebugBuild()
		debugQuery(debug)
	}

	err = wrk.db.QueryRow(ctx, stmt, workspaceId).Scan(
		&setting.WorkspaceId, &setting.ServerId, &setting.ServerToken, &setting.IsEnabled, &setting.Email,
		&setting.Domain,
		&setting.HasError, &inboundEmail,
		&setting.HasForwardingEnabled,
		&setting.HasDNS, &setting.IsDNSVerified, &dnsVerifiedAt, &dnsDomainId,
		&dkimHost, &dkimTextValue, &dkimUpdateStatus,
		&returnPathDomain, &returnPathDomainCname, &setting.ReturnPathDomainVerified,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return models.PostmarkMailServerSetting{}, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to query", slog.Any("err", err))
		return models.PostmarkMailServerSetting{}, ErrQuery
	}

	if inboundEmail.Valid {
		setting.InboundEmail = &inboundEmail.String
	}
	if dnsVerifiedAt.Valid {
		setting.DNSVerifiedAt = &dnsVerifiedAt.Time
	}
	if dnsDomainId.Valid {
		setting.DNSDomainId = &dnsDomainId.Int64
	}
	if dkimHost.Valid {
		setting.DKIMHost = &dkimHost.String
	}
	if dkimTextValue.Valid {
		setting.DKIMTextValue = &dkimTextValue.String
	}
	if dkimUpdateStatus.Valid {
		setting.DKIMUpdateStatus = &dkimUpdateStatus.String
	}
	if returnPathDomain.Valid {
		setting.ReturnPathDomain = &returnPathDomain.String
	}
	if returnPathDomainCname.Valid {
		setting.ReturnPathDomainCNAME = &returnPathDomainCname.String
	}
	return setting, nil
}
