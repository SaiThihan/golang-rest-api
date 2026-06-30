package store

import "database/sql"

type password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	CreatedAt    string `json:"created_at"`
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
	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: "",
	}

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
