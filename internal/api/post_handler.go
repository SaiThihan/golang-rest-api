package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SaiThihan/go-basic/internal/store"
	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	postStore store.PostStore
}

func NewPostHandler(ps store.PostStore) *PostHandler {
	return &PostHandler{
		postStore: ps,
	}
}

// posts/:id
func (ph *PostHandler) HandleGetPostById(w http.ResponseWriter, r *http.Request) {
	paramPostId := chi.URLParam(r, "id")
	if paramPostId == "" {
		http.NotFound(w, r)
		return
	}

	postId, err := strconv.ParseInt(paramPostId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Post ID: %d\n", postId)
}

func (ph *PostHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	var post store.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid JSON Request", http.StatusBadRequest)
		return
	}

	createdPost, err := ph.postStore.CreatePost(&post)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPost)
}
