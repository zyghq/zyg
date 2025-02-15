package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/cristalhq/builq"
	"github.com/jackc/pgx/v5"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/models"
)

func workosUserCols() builq.Columns {
	return builq.Columns{
		"user_id",
		"email",
		"first_name",
		"last_name",
		"email_verified",
		"profile_picture_url",
		"created_at",
		"updated_at",
	}
}

func (u *UserDB) SaveWorkOSUser(ctx context.Context, user *models.WorkOSUser) (*models.WorkOSUser, error) {
	var profilePictureUrl sql.NullString

	if user.ProfilePictureURL != nil {
		profilePictureUrl = sql.NullString{
			String: *user.ProfilePictureURL,
			Valid:  true,
		}
	}

	cols := workosUserCols()
	insertB := builq.Builder{}
	insertParams := []any{
		user.UserID, user.Email, user.FirstName, user.LastName, user.EmailVerified, profilePictureUrl,
		user.CreatedAt, user.UpdatedAt,
	}

	insertB.Addf("INSERT INTO workos_user (%s)", cols)
	insertB.Addf("VALUES (%$, %$, %$, %$, %$, %$, %$, %$)", insertParams...)
	insertB.Addf("ON CONFLICT (user_id) DO UPDATE SET")
	insertB.Addf("email = EXCLUDED.email,")
	insertB.Addf("first_name = EXCLUDED.first_name,")
	insertB.Addf("last_name = EXCLUDED.last_name,")
	insertB.Addf("email_verified = EXCLUDED.email_verified,")
	insertB.Addf("profile_picture_url = EXCLUDED.profile_picture_url,")
	insertB.Addf("created_at = EXCLUDED.created_at,")
	insertB.Addf("updated_at = EXCLUDED.updated_at")
	insertB.Addf("RETURNING %s", cols)

	insertQuery, _, err := insertB.Build()
	if err != nil {
		slog.Error("failed to build upsert query", slog.Any("err", err))
		return nil, ErrQuery
	}

	if zyg.DBQueryDebug() {
		debug := insertB.DebugBuild()
		debugQuery(debug)
	}

	err = u.db.QueryRow(ctx, insertQuery, insertParams...).Scan(
		&user.UserID, &user.Email, &user.FirstName, &user.LastName, &user.EmailVerified, &user.ProfilePictureURL,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Error("no rows returned", slog.Any("err", err))
		return nil, ErrEmpty
	}
	if err != nil {
		slog.Error("failed to insert query", slog.Any("err", err))
		return nil, ErrQuery
	}

	if profilePictureUrl.Valid {
		user.ProfilePictureURL = &profilePictureUrl.String
	}
	return user, nil
}
