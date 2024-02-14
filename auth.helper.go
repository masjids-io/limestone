package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mnadev/limestone/auth"
	"github.com/mnadev/limestone/storage"
	"github.com/oov/gothic"
	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

// Routes for authentication
func (s *App) HandleGoogleOauthRoute(w http.ResponseWriter, r *http.Request) {
	err := gothic.BeginAuth("google", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *App) HandleGoogleOauthCallbackRoute(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteAuth("google", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var dbUser storage.User
	err = s.DB.Where("email = ?", user.Email).First(&dbUser).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		// uuid generation
		uuid, err := uuid.NewUUID()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// create an entry for a user
		newUser := storage.User{
			ID:             uuid,
			Email:          user.Email,
			Username:       user.NickName,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			HashedPassword: "social_auth",
			IsVerified:     true,
			PhoneNumber:    "",
			Gender:         storage.MALE, //defaulting to male
		}
		//create a db entry
		err = s.DB.Create(newUser).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//uuid taken after creating entry here (SIMULATES REGISTER + LOGIN)
		access, _ := auth.CreateJWT(uuid.String(), user.Email, auth.TokenType(0))  //access token
		refresh, _ := auth.CreateJWT(uuid.String(), user.Email, auth.TokenType(1)) //refresh token
		response := map[string]string{
			"accessToken":  access,
			"refreshToken": refresh,
		}
		// Set the Content-Type header and response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return

	}
	//Create JWT tokens and return in response from fetched user (SIMULATES LOGIN)
	access, _ := auth.CreateJWT(dbUser.ID.String(), user.Email, auth.TokenType(0))  //access token
	refresh, _ := auth.CreateJWT(dbUser.ID.String(), user.Email, auth.TokenType(1)) //refresh token
	// Prepare the response
	response := map[string]string{
		"accessToken":  access,
		"refreshToken": refresh,
	}
	// Set the Content-Type header and response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
