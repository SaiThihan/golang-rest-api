package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/SaiThihan/go-basic/internal/store"
	"github.com/SaiThihan/go-basic/internal/utils"
)

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) < 4 {
		return errors.New("username must be at least 4 characters long")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9].[a-zA-Z0-9\._%\+\-]{0,63}@[a-zA-Z0-9\.\-]+\.[a-zA-Z]{2,30}$`)

	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

func (uh *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Payload{"error": "Invalid request payload"})
		return
	}

	err = uh.validateRegisterRequest(&req)
	if err != nil {
		uh.logger.Printf("Invalid register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Payload{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	err = user.PasswordHash.HashPassword(req.Password)
	if err != nil {
		uh.logger.Printf("Error hashing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Payload{"error": "Internal server error"})
		return
	}

	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("Error creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Payload{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Payload{"user": user})
}
