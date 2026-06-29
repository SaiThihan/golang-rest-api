package store

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type PostStore interface {
	CreatePost(*Post) (*Post, error)
	GetPostById(id int64) (*Post, error)
}

type PostgresPostStore struct {
	db *sql.DB
}

func NewPostgresPostStore(db *sql.DB) *PostgresPostStore {
	return &PostgresPostStore{db: db}
}

func (pg *PostgresPostStore) CreatePost(post *Post) (*Post, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	// Ensure that the transaction is rolled back if an error occurs
	defer tx.Rollback()

	query := `INSERT INTO posts (title, content) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(query, post.Title, post.Content).Scan(&post.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pg *PostgresPostStore) GetPostById(id int64) (*Post, error) {
	post := &Post{}
	return post, nil
}
