package main

import (
	"net/http"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var login = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	if web.Authorized(w, r, []string{"USER", "ADMIN"}) {
		http.Redirect(w, r, "/my-list", 303)
		return
	}
	tmpl.Render(w, r, "login.tmpl", nil)
	return
}}

var register = web.Route{"POST", "/register", func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var user User
	if errs, ok := FormToStruct(&user, r.Form, "register"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/", "Error registering")
		return
	}
	// check for uniqueness here
	var users []User
	db.TestQuery("user", &user, adb.Eq("email", user.Email), adb.Ne("id", `"`+user.Id+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/", "Error registering. Email is already in use.")
	}
	user.Active = true
	user.Id = genId()
	user.Role = "USER"
	if !db.Add("user", user.Id, user) {
		web.SetErrorRedirect(w, r, "/", "Error registering. Please try again")
		return
	}
	// success login user in
	sess := web.Login(w, r, user.Role)
	sess["id"] = user.Id
	sess["email"] = user.Email
	web.PutMultiSess(w, r, sess)
	web.SetSuccessRedirect(w, r, "/my-list", "welcome "+user.FirstName)
	return
}}

var loginPost = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var user User
	if errs, ok := FormToStruct(&user, r.Form, "login"); !ok {
		web.SetFormErrors(w, errs)
		web.SetErrorRedirect(w, r, "/", "Error logging in")
		return
	}
	if !db.Auth("user", user.Email, user.Password, &user) {
		web.SetErrorRedirect(w, r, "/", "Incorrect email or password")
		return
	}
	sess := web.Login(w, r, user.Role)
	sess["id"] = user.Id
	sess["email"] = user.Email
	web.PutMultiSess(w, r, sess)
	web.SetSuccessRedirect(w, r, "/my-list", "welcome "+user.FirstName)
	return
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w)
	web.SetSuccessRedirect(w, r, "/", "Successfully logged out")
	return
}}
