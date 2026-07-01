package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SaiThihan/go-basic/internal/store"
	"github.com/SaiThihan/go-basic/internal/tokens"
	"github.com/SaiThihan/go-basic/internal/utils"
)

type TokenHandler struct {
	userStore  store.UserStore
	tokenStore store.TokenStore
	logger     *log.Logger
}

func NewTokenHandler(userStore store.UserStore, tokenStore store.TokenStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		userStore:  userStore,
		tokenStore: tokenStore,
		logger:     logger,
	}
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (th *TokenHandler) HandleTokenCreate(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		th.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Payload{"error": "Invalid request payload"})
		return
	}

	user, err := th.userStore.GetUserByUsername(req.Username)
	if err != nil {
		th.logger.Printf("Error fetching user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Payload{"error": "Internal server error"})
		return
	}
	if user == nil {
		th.logger.Printf("User not found for username: %q", req.Username)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Payload{"error": "Invalid credentials"})
		return
	}

	isMatch, err := user.PasswordHash.ComparePassword(req.Password)
	if err != nil {
		th.logger.Printf("Error comparing passwords: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Payload{"error": "Internal server error"})
		return
	}
	if !isMatch {
		th.logger.Printf("Invalid password for user: %s", req.Username)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Payload{"error": "Invalid credentials"})
		return
	}

	token, err := th.tokenStore.CreateToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		th.logger.Printf("Error creating token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Payload{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Payload{"token": token})

}
