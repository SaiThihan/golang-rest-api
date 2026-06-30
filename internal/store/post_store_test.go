package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	connectionString := "host=localhost user=test_user password=test_pass dbname=test_postgres port=5435 sslmode=disable"
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestCreatePost(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresPostStore(db)

	tests := []struct {
		name    string
		post    *Post
		wantErr bool
	}{
		{
			name: "Create a new valid post",
			post: &Post{
				Title:   "Test Post",
				Content: "This is a test post.",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdPost, err := store.CreatePost(tt.post)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.post.Title, createdPost.Title)
			assert.Equal(t, tt.post.Content, createdPost.Content)

			fetchedPost, err := store.GetPostById(createdPost.ID)

			require.NoError(t, err)
			assert.Equal(t, createdPost.Title, fetchedPost.Title)
			assert.Equal(t, createdPost.Content, fetchedPost.Content)
		})
	}
}
