package repositoryuser

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/internal/repository"
)

type userRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db,
	}
}

func (repo *userRepository) Insert(user *data.User) error {
	query := `
		INSERT INTO users (username, role, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, role, created_at`

	args := []interface{}{user.Username, user.Role, user.Password.Hash}
	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Role, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return data.ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (repo *userRepository) Get(username string) (*data.User, error) {

	query := `SELECT id, username, deposit, password_hash, role, created_at
			  FROM users
			  WHERE username = $1`

	var user data.User

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	err := repo.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Deposit,
		&user.Password.Hash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (repo *userRepository) Update(user *data.User) error {
	query := `
		UPDATE users
		SET username = $1, password_hash = $2, deposit = $3
		WHERE id = $4 RETURNING username, deposit`

	args := []interface{}{user.Username, user.Password.Hash, user.Deposit, user.ID}

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&user.Username, &user.Deposit)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) GetForToken(tokenPlainText, scope string) (*data.User, error) {

	hash := sha256.Sum256([]byte(tokenPlainText))

	query := `
			SELECT users.id, users.created_at, users.username, users.role, 
			users.password_hash, users.deposit
			FROM users
			INNER JOIN tokens
			ON users.id = tokens.user_id
			WHERE tokens.hash = $1
			AND tokens.scope = $2
			AND tokens.expiry > $3`

	args := []interface{}{hash[:], scope, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	var user data.User

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Role,
		&user.Password.Hash,
		&user.Deposit,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

func (repo *userRepository) Delete(id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return data.ErrRecordNotFound
	}

	return nil
}
