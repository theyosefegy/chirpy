package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
	myutils "github.com/theyosefegy/chriby/util"
)

var userCounter = 1 

type User struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Expires_in_seconds float64 `json:"expiresinseconds"` 
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	UserToken string `json:"token"`
}

var users = []User{}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		myutils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Check if the email already exists
	for _, u := range users {
		if u.Email == user.Email {
			myutils.RespondWithError(w, http.StatusBadRequest, "Email already exists")
			return
		}
	}

	user.ID = userCounter
	userCounter++

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		myutils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user.Password = string(hashedPassword)

	// توليد التوكن
	token, err := generateToken(user)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		myutils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	users = append(users, user)

	response := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		UserToken: token,
	}

	myutils.RespondWithJSON(w, 201, response)
}
func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq User = User{
		Expires_in_seconds: 24,
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginReq)

	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		myutils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Find user by email
	var user *User
	for _, u := range users {
		if u.Email == loginReq.Email {
			user = &u
			break
		}
	}

	if user == nil {
		myutils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		myutils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := generateToken(*user)

	if err != nil {
		myutils.RespondWithError(w, http.StatusForbidden, "Something went wrong while making the user's token.")
	}

	response := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		UserToken: token,
	}

	myutils.RespondWithJSON(w, 200, response)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if users == nil {
		myutils.RespondWithError(w, 404, "There's no users to show, yet")
		return
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:    user.ID,
			Email: user.Email,
		}
	}

	myutils.RespondWithJSON(w, 200, userResponses)
}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		myutils.RespondWithError(w, 404, "Invalid id")
		return
	}

	for _, user := range users {
		if user.ID == id {
			response := UserResponse{
				ID:    user.ID,
				Email: user.Email,
			}
			myutils.RespondWithJSON(w, 201, response)
			return
		}
	}

	myutils.RespondWithError(w, 400, "No user with the specific ID.")
}


func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extract and validate JWT
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		myutils.RespondWithError(w, http.StatusUnauthorized, "Authorization header missing")
		return
	}

	// Remove "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		myutils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	// Extract user ID from token claims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.Subject == "" {
		myutils.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	userID := claims.Subject

	// 2. Find the user
	var user *User
	for i, u := range users {
		if string(u.ID) == userID {
			user = &users[i]
			break
		}
	}

	if user == nil {
		myutils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// 3. Update the user
	var updateData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&updateData)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		myutils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if updateData.Email != "" {
		user.Email = updateData.Email
	}

	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %s", err)
			myutils.RespondWithError(w, http.StatusInternalServerError, "Error updating password")
			return
		}
		user.Password = string(hashedPassword)
	}

	// 4. Return updated user data (excluding the password)
	response := UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
	myutils.RespondWithJSON(w, http.StatusOK, response)
}