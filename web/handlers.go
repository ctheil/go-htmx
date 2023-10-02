package web

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	todos "github.com/ctheil/go-htmx"
	"github.com/ctheil/go-htmx/tools"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func NewHandler(store todos.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	h.Use(middleware.Logger)
	h.Use(h.Authorize)

	h.Get("/", h.Home())
	h.Get("/login", h.Login())
	h.Post("/login", h.PostLogin())
	h.Get("/signup", h.Signup())
	h.Post("/signup", h.PostSignup())
	//
	// HTMX PARTIAL ROUTES //
	//
	h.Route("/todos", func(r chi.Router) {
		r.Get("/", h.TodosList())
		r.Post("/", h.PostTodo())
		r.Get("/create", h.CreateTodo())
	})

	return h
}

type Handler struct {
	*chi.Mux

	store todos.Store
}

func (h *Handler) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			ctx := context.WithValue(r.Context(), "user", nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		sessionToken := uuid.MustParse(c.Value)

		session, err := h.store.GetSession(sessionToken)
		if err != nil {
			// e := tools.ClientError{
			// 	Heading: "Not signed in",
			// 	Message: "Please sign in",
			// 	Href:    "/login",
			// }
			// BuildError(w, e)
			ctx := context.WithValue(r.Context(), "user", nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return

		}
		u, err := h.store.GetUserByEmail(session.Email)
		if err != nil {
			// e := tools.ClientError{
			// 	Heading: "Not signed in",
			// 	Message: "Please sign in",
			// 	Href:    "/login",
			// }
			// BuildError(w, e)
			ctx := context.WithValue(r.Context(), "user", u.ID.String())
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx := context.WithValue(r.Context(), "user", u.ID.String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) PostTodo() http.HandlerFunc {
	fmt.Println("POST TODO")
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			e := tools.ClientError{
				Heading: "Invalid input",
			}
			BuildError(w, e)
			return
		}

		dueStr := r.FormValue("due")
		dueTime, err := time.Parse("2006-01-02", dueStr)
		if err != nil {
			fmt.Println("Invalid date formatting", err)
			e := tools.ClientError{
				Heading: "Invalid date formatting",
				Message: err.Error(),
			}
			BuildError(w, e)
			return
		}

		id := uuid.New()
		uid := uuid.MustParse(r.Context().Value("user").(string))
		t := todos.Todo{
			ID:          id,
			UserID:      uid,
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Complete:    false,
			Due:         dueTime,
		}
		fmt.Println("TODO: ", t)
		if err := h.store.CreateTodo(&t); err != nil {
			fmt.Println(err)
			e := tools.ClientError{
				Heading: "Error creating todo",
			}
			BuildError(w, e)
			return
		}

	}
}

func (h *Handler) CreateTodo() http.HandlerFunc {
	fmt.Println("Create Todo!")
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/create-todo.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "create-todo", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

func (h *Handler) TodosList() http.HandlerFunc {

	fmt.Println("TODOS!")

	return func(w http.ResponseWriter, r *http.Request) {
		uidStr := r.Context().Value("user").(string)
		uid := uuid.MustParse(uidStr)
		u, err := h.store.GetUserById(uid)
		if err != nil {
			fmt.Println("no user")
		}

		tt, err := h.store.GetTodos(u.ID)
		if err != nil {
			e := tools.ClientError{
				Heading: "No todos!",
				Message: "Create todos?",
				Href:    "/todos/create",
			}
			BuildError(w, e)
			return
		}

		// return tt
		tmpl, err := template.ParseFiles("templates/home.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "todo", tt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *Handler) Home() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		uidStr := r.Context().Value("user").(string)
		uid := uuid.MustParse(uidStr)
		u, err := h.store.GetUserById(uid)
		if err != nil {
			fmt.Println("No user")
		}
		fmt.Println("USER:", u)
		tt, err := h.store.GetTodos(u.ID)
		if err != nil {
			e := tools.ClientError{
				Heading: "No todos!",
				Message: "Create todos?",
				Href:    "/todos/create",
			}
			BuildError(w, e)
		}
		u.Todos = tt

		fmt.Println("HOME", u)
		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, u); err != nil {
			fmt.Println("Error executing template")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}
