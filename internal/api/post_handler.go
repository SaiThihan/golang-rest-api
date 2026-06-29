package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
}

func NewPostHandler() *PostHandler {
	return &PostHandler{}
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
	fmt.Fprintln(w, "Post created successfully")
}
