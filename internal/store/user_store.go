package store

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) HashPassword(plainText string) error {
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)

	if err != nil {
		return err
	}

	p.plaintext = &plainText
	p.hash = hashpassword

	return nil
}

func (p *password) ComparePassword(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type User struct {
	ID           int      `json:"id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	PasswordHash password `json:"-"`
	CreatedAt    string   `json:"created_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `INSERT INTO users (username, email, password_hash ) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{}

	query := `SELECT id,username,email, created_at FROM users WHERE username = $1`
	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3`

	result, err := s.db.Exec(query, user.Username, user.Email, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
