package main

import (
	"log"
	"net/http"
	"net/url"
)

func (*apiCfg) register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		w.WriteHeader(500)
		renderRegisterFail(w, r, "Could not processs registration. Please try again later.")
		return
	}

	valid, errMsg := validate(r.Form)
	if !valid {
		renderRegisterFail(w, r, errMsg)
		return
	}
	// put in database
	w.WriteHeader(201)
	renderRegisterSuccess(w, r)
	
}

func validate(data url.Values) (bool, string) {
	var msg string
	var res bool = true

	username := data["username"][0]
	password := data["password"][0]

	usernameLength := len(username)
	passwordLength := len(password)
	log.Println(usernameLength, passwordLength)
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