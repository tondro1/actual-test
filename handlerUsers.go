package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/tondro1/actual-test/internal/auth"
	"github.com/tondro1/actual-test/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (api *apiCfg) register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		w.WriteHeader(500)
		renderRegisterFail(w, "Could not processs registration. Please try again later.")
		return
	}

	username := r.Form["username"][0]
	password := r.Form["password"][0]

	valid, errMsg := validate(username, password)
	if !valid {
		renderRegisterFail(w, errMsg)
		return
	}
	// put in database
	errMsg = "Could not register new user."
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		renderRegisterFail(w, errMsg)
		log.Println(err)
		return
	}
	
	_, err = api.db.CreateUser(r.Context(), database.CreateUserParams{
		ID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		CreatedAt: pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		Username: username,
		Password: hash,
	})
	
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusBadRequest)
			renderRegisterFail(w, errMsg + " Username has been taken. Please choose a unique username.")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		renderRegisterFail(w, errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
	renderRegisterSuccess(w)
}

func (api *apiCfg) login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		w.WriteHeader(http.StatusInternalServerError)
		renderRegisterFail(w, "Could not processs login. Please try again later.")
		return
	}

	username := r.Form["username"][0]
	password := r.Form["password"][0]

	// validate
	user, err := api.db.GetUser(r.Context(), username)
	if err != nil {
		log.Println("Error getting user from db:", err)
		w.WriteHeader(http.StatusInternalServerError)
		renderRegisterFail(w, "Could not processs login. Please try again later.")
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		log.Println("Error password incorrect:", err)
		w.WriteHeader(http.StatusInternalServerError)
		renderRegisterFail(w, "Password is incorrect, please try again.")
		return
	}

	tokenString, err := auth.NewAccessToken(user)
	if err != nil {
		log.Println(err)
		return
	}

	jwt := http.Cookie {
		Name: "accessToken",
		Path: "/",
		Value: tokenString,
		MaxAge: 36000,
		SameSite: http.SameSiteStrictMode,
		Secure: true,
		HttpOnly: true,
	}
	http.SetCookie(w, &jwt)
	w.Header().Add("HX-Location", "/")
	w.WriteHeader(http.StatusAccepted)
}

func validate(username string, password string) (bool, string) {
	var msg string
	var res bool = true

	usernameLength := len(username)
	passwordLength := len(password)
	if usernameLength < 8 || usernameLength > 40 {
		res = false
		msg += "Username must be between 8 and 40 characters long.\n"
	}

	if passwordLength < 8 {
		res = false
		msg += "Password must be longer than 8 characters."
	}

	return res, msg
}
