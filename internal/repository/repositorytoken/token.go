package repositorytoken

import (
	"context"
	"database/sql"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/internal/repository"
)

type tokenRepository struct {
	*sql.DB
}

func NewTokenRepository(db *sql.DB) repository.TokenRepository {
	return &tokenRepository{db}
}

func (repo *tokenRepository) Create(token *data.Token) error {

	query := `
			INSERT INTO tokens (hash, user_id, expiry, scope)
			VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserId, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, query, args...)

	return err
}

func (repo *tokenRepository) DeleteAllForUserByScope(scope string, userID int64) error {

	query := `
			DELETE FROM tokens
			WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, query, scope, userID)

	return err
}
