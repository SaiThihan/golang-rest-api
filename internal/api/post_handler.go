package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SaiThihan/go-basic/internal/store"
	"github.com/SaiThihan/go-basic/internal/utils"
)

type PostHandler struct {
	postStore store.PostStore
	logger    *log.Logger
}

func NewPostHandler(ps store.PostStore, logger *log.Logger) *PostHandler {
	return &PostHandler{
		postStore: ps,
		logger:    logger,
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
		ph.logger.Printf("Error creating post: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Payload{"error": "Failed to create post"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Payload{"createdPost": createdPost})
}

func (ph *PostHandler) HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := ph.postStore.GetPosts()
	if err != nil {
		ph.logger.Printf("Error retrieving posts: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Payload{"error": "Failed to retrieve posts"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Payload{"posts": posts})
}

// posts/:id
func (ph *PostHandler) HandleGetPostById(w http.ResponseWriter, r *http.Request) {
	postId, err := utils.RetrieveIDFromRequest(r)
	if err != nil {
		return
	}

	post, err := ph.postStore.GetPostById(postId)
	if err != nil {
		ph.logger.Printf("Error retrieving post by ID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Payload{"error": "Failed to retrieve post"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if post == nil {
		http.NotFound(w, r)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Payload{"post": post})
}

func (ph *PostHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	postId, err := utils.RetrieveIDFromRequest(r)

	if err != nil {
		return
	}

	err = ph.postStore.DeletePost(postId)

	if err != nil {
		ph.logger.Printf("Error deleting post: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Payload{"error": "Failed to delete post"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Payload{"message": "Post deleted successfully"})
}

func (ph *PostHandler) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	postId, err := utils.RetrieveIDFromRequest(r)

	if err != nil {
		return
	}

	var post store.Post

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid JSON Request", http.StatusBadRequest)
		return
	}

	post.ID = postId

	updatedPost, err := ph.postStore.UpdatePost(&post)
	if err != nil {
		ph.logger.Printf("Error updating post: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Payload{"error": "Failed to update post"})
		return
	}

	if updatedPost == nil {
		http.NotFound(w, r)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Payload{"post": updatedPost})
}
