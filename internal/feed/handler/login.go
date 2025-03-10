package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/mail"
	"news/platform/web"
)

type LoginPageData struct {
	Email        string
	Error        string
	Token        string
	JustLoggedIn bool
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	return renderLoginPage(w, LoginPageData{}, h.Cache)
}

func (h *Handler) DoLogin(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	email := r.FormValue("email")

	if !isValidEmail(email) {
		if len(email) > 40 { //is a token, TODO validate in the DB
			return renderLoginPage(w, LoginPageData{
				Email: email,
				Error: "Error processing login, please try again",
			}, h.Cache)
		}

		return renderLoginPage(w, LoginPageData{
			Email: email,
			Error: "Please enter a valid email address",
		}, false)
	}

	token, err := generateToken()
	if err != nil {
		return renderLoginPage(w, LoginPageData{
			Email: email,
			Error: "Error processing login, please try again",
		}, h.Cache)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400 * 365, // 1 year
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   h.SecureCookie,
	})

	// Redirect to home with the token
	http.Redirect(w, r, "/?loggedIn=true", http.StatusSeeOther)
	return nil
}

func (h *Handler) Token(w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie("auth_token")
	if err != nil {
		return err
	}

	if len(c.Value) != 44 {
		return web.TemplateRender(w, "home.page.tmpl", &web.TemplateData{}, h.Cache)
	}

	return web.TemplateRender(w, "token.page.tmpl", &web.TemplateData{Token: c.Value, JustLoggedIn: true}, h.Cache)
}

func renderLoginPage(w http.ResponseWriter, data LoginPageData, cache bool) error {
	err := web.TemplateRender(w, "login.page.tmpl", &web.TemplateData{
		Email:        data.Email,
		LoginErr:     data.Error,
		Token:        data.Token,
		JustLoggedIn: data.JustLoggedIn,
	}, cache)

	if err != nil {
		return err
	}

	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func generateToken() (string, error) {
	// Generate 32 bytes of random data
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Convert to base64 for easier handling
	return base64.URLEncoding.EncodeToString(b), nil
}
