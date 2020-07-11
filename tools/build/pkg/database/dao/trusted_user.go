package dao

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/xerrors"

	"go.f110.dev/mono/tools/build/pkg/database"
)

type TrustedUser struct {
	conn *sql.DB
}

func NewTrustedUser(conn *sql.DB) *TrustedUser {
	return &TrustedUser{conn: conn}
}

func (t *TrustedUser) SelectByGithubId(ctx context.Context, githubId int64) (*database.TrustedUser, error) {
	row := t.conn.QueryRowContext(ctx, "SELECT `id`, `created_at`, `updated_at` FROM `trusted_user` WHERE `github_id` = ?", githubId)

	trustedUser := &database.TrustedUser{GithubId: githubId}
	if err := row.Scan(&trustedUser.Id, &trustedUser.CreatedAt, &trustedUser.UpdatedAt); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return trustedUser, nil
}

func (t *TrustedUser) Create(ctx context.Context, user *database.TrustedUser) (*database.TrustedUser, error) {
	res, err := t.conn.ExecContext(ctx, "INSERT INTO `trusted_user` (`github_id`, `created_at`) VALUES (?, ?)", user.GithubId, time.Now())
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return nil, sql.ErrNoRows
	}

	user = user.Copy()
	insertedId, err := res.LastInsertId()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	user.Id = int32(insertedId)

	user.ResetMark()
	return user, nil
}