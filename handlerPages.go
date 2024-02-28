package main

import (
	"html/template"
	"net/http"

	"github.com/tondro1/actual-test/internal/auth"
)

func renderIndex(w http.ResponseWriter, r *http.Request, uclaims *auth.UserClaims) {
	var tmpl *template.Template
	if uclaims.UserId != "" {
		//logged in
		tmpl = template.Must(template.ParseFiles("./public/root.html", "./public/loggedin/index.html", "./public/loggedin/navbar.html"))
	} else {
		// not logged in
		tmpl = template.Must(template.ParseFiles("./public/root.html", "./public/index.html", "./public/navbar.html"))
	}

	if len(r.Header["Hx-Request"]) == 1 {
		tmpl.ExecuteTemplate(w, "body", nil)
	} else {
		tmpl.ExecuteTemplate(w, "root", nil)
	}
	
}

func renderLogin(w http.ResponseWriter, r *http.Request, uclaims *auth.UserClaims)  {
	var tmpl *template.Template
	if uclaims.UserId != "" {
		renderIndex(w, r, uclaims)
		return	
	} else {
		tmpl = template.Must(template.ParseFiles("./public/root.html", "./public/login.html", "./public/navbar.html"))
	}
	if len(r.Header["Hx-Request"]) == 1 {
		tmpl.ExecuteTemplate(w, "body", nil)
	} else {
		tmpl.ExecuteTemplate(w, "root", nil)
	}
}

func renderRegister(w http.ResponseWriter, r *http.Request, uclaims *auth.UserClaims) {
	var tmpl *template.Template
	if uclaims.UserId != "" {
		renderIndex(w, r, uclaims)
		return
	} else {
		tmpl = template.Must(template.ParseFiles("./public/root.html", "./public/register.html", "./public/navbar.html"))
	}
	
	if len(r.Header["Hx-Request"]) == 1 {
		tmpl.ExecuteTemplate(w, "body", nil)
	} else {
		tmpl.ExecuteTemplate(w, "root", nil)
	}
}

func renderRegisterSuccess(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("./public/register.success.html"))
	tmpl.Execute(w, nil)
}

func renderRegisterFail(w http.ResponseWriter, msg string) {
	tmpl := template.Must(template.ParseFiles("./public/register.fail.html"))
	type Data struct {
		Message string
	}

	tmpl.Execute(w, Data{Message: msg})
}
