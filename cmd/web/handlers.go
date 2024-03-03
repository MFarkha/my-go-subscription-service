package main

import "net/http"

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	// set up homepage
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	// set up homepage
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	//
	_ = app.Session.RenewToken(r.Context())

	// parse post
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println("error parsing post data:", err)
	}

	// get email & password
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	validPassword, err := user.PasswordMatches(password)
	if err != nil || !validPassword {
		app.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	// ok log user in
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)
	app.Session.Put(r.Context(), "flash", "Successful login!")

	// redirect the user
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) LogOut(w http.ResponseWriter, r *http.Request) {
	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	// create user

	// send an activation email concurrently

	// subscribe the user to an account
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate the url

	// generate an invoice (send email)

	// send an email with attachments

	// send an email with the invoice attached
}
