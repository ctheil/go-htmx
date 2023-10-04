package web

import (
	"fmt"
	"html/template"

	"net/http"

	todos "github.com/ctheil/go-htmx"
	"github.com/ctheil/go-htmx/auth"
	"github.com/ctheil/go-htmx/tools"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func BuildError(w http.ResponseWriter, e tools.ClientError) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "error", e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) Login() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/login.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}
func (h *Handler) PostLogin() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		u, err := h.store.GetUserByEmail(r.FormValue("email"))
		if err != nil {
			fmt.Println("User not found")
			clientError := tools.ClientError{
				Heading: "No User Found",
				Message: `Signup instead?`,
				Href:    "/signup",
			}
			BuildError(w, clientError)

			return
		}
		fmt.Println("User found")

		isPasswordValid := auth.CheckPasswordHash(r.FormValue("password"), u.Password)
		fmt.Println("Password is valid: ", isPasswordValid)

		if !isPasswordValid {
			e := tools.ClientError{
				Heading: "Invalid Password",
				Message: "Please, try again.",
			}
			BuildError(w, e)

			return
		}
		fmt.Println("Everything checks out")
		// generate cookie
		s, err := h.store.StartSession(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}
		fmt.Println(s)
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   s.Token.String(),
			Expires: s.Expiry,
		})

		w.Header().Set("HX-Redirect", "/")
		w.Header().Set("HX-Location", "/")
		w.WriteHeader(http.StatusOK)

	}

}
func (h *Handler) Signup() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/signup.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}
func (h *Handler) PostSignup() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}
		if r.FormValue("password") != r.FormValue("confirm-password") {
			// return error
			clientError := tools.ClientError{
				Heading: "Passwords do not match",
				Message: "",
				Href:    "",
			}
			BuildError(w, clientError)
			return
		}
		_, err := h.store.GetUserByEmail(r.FormValue("email"))

		if err == nil {
			fmt.Println(err)

			e := tools.ClientError{
				Heading: "User already exists",
				Message: "Login instead?",
				Href:    "/login",
			}
			BuildError(w, e)
			return
		}
		id := uuid.New()
		hashedPassword, err := auth.HashPassword(r.FormValue("password"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u := todos.User{
			Email:    r.FormValue("email"),
			Password: hashedPassword,
			ID:       id,
			Name:     r.FormValue("name"),
		}
		fmt.Println(u)
		if err := h.store.CreateUser(&u); err != nil {
			fmt.Println(err)

			e := tools.ClientError{
				Heading: "Error creating user: ",
				Message: err.Error(),
			}
			BuildError(w, e)
			return
		}
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)

	}

}
