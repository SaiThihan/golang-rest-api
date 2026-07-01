package store

import (
	"database/sql"
	"time"

	"github.com/SaiThihan/go-basic/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	CreateToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	SaveToken(token *tokens.Token) error
	DeleteTokenwithUserId(userID int, scope string) error
}

func (ts *PostgresTokenStore) CreateToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = ts.SaveToken(token)

	return token, nil
}

func (ts *PostgresTokenStore) SaveToken(token *tokens.Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`
	_, err := ts.db.Exec(query, token.Hash, token.UserId, token.Expiry, token.Scope)
	return err
}

func (ts *PostgresTokenStore) DeleteTokenwithUserId(userID int, scope string) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`
	_, err := ts.db.Exec(query, userID, scope)
	return err
}
