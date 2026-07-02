package store

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PostStore interface {
	CreatePost(*Post) (*Post, error)
	GetPosts() ([]Post, error)
	GetPostById(id int64) (*Post, error)
	UpdatePost(post *Post) (*Post, error)
	DeletePost(id int64) error
	GetPostOwner(id int64) (int, error)
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

	query := `INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(query, post.Title, post.Content, post.UserID).Scan(&post.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pg *PostgresPostStore) GetPosts() ([]Post, error) {
	query := `SELECT id, title, content, created_at FROM posts ORDER BY created_at DESC`

	rows, err := pg.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)

		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (pg *PostgresPostStore) GetPostById(id int64) (*Post, error) {
	post := &Post{}
	query := `SELECT id, title, content, created_at FROM posts WHERE id = $1`

	err := pg.db.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pg *PostgresPostStore) DeletePost(id int64) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `DELETE FROM posts WHERE id = $1`
	result, err := tx.Exec(query, id)
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

	return tx.Commit()
}

func (pg *PostgresPostStore) UpdatePost(post *Post) (*Post, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `UPDATE posts SET title = $1, content = $2 WHERE id= $3 RETURNING id, title, content, created_at`

	result, err := tx.Exec(query, post.Title, post.Content, post.ID)

	if err != nil {
		return nil, err
	}

	rowsEffected, err := result.RowsAffected()

	if err != nil {
		return nil, err
	}

	if rowsEffected == 0 {
		return nil, nil
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pg *PostgresPostStore) GetPostOwner(postID int64) (int, error) {
	var userID int
	query := `SELECT user_id FROM posts WHERE id = $1`

	if err := pg.db.QueryRow(query, postID).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}
