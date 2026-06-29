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

func (ph *PostHandler) HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := ph.postStore.GetPosts()
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
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

	post, err := ph.postStore.GetPostById(postId)
	if err != nil {
		http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if post == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (ph *PostHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
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

	err = ph.postStore.DeletePost(postId)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(nil)
}
