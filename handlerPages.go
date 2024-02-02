package main

import (
	"html/template"
	"net/http"
)

func renderIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./public/root.html", "./public/index.html", "./public/navbar.html"))
	if len(r.Header["Hx-Request"]) == 1 {
		tmpl.ExecuteTemplate(w, "body", nil)
	} else {
		tmpl.ExecuteTemplate(w, "root", nil)
	}
}

func renderLogin(w http.ResponseWriter, r *http.Request)  {
	tmpl := template.Must(template.ParseFiles("./public/root.html", "./public/login.html", "./public/navbar.html"))
	if len(r.Header["Hx-Request"]) == 1 {
		tmpl.ExecuteTemplate(w, "body", nil)
	} else {
		tmpl.ExecuteTemplate(w, "root", nil)
	}
}

func renderRegister(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./public/root.html", "./public/register.html", "./public/navbar.html"))
	if len(r.Header["Hx-Request"]) == 1 {
		tmpl.ExecuteTemplate(w, "body", nil)
	} else {
		tmpl.ExecuteTemplate(w, "root", nil)
	}
}

func renderRegisterSuccess(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./public/register.success.html"))
	tmpl.Execute(w, nil)
}

func renderRegisterFail(w http.ResponseWriter, r *http.Request, msg string) {
	tmpl := template.Must(template.ParseFiles("./public/register.fail.html"))
	type Data struct {
		Message string
	}

	tmpl.Execute(w, Data{Message: msg})
}