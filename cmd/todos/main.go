package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	// locals
	"github.com/ctheil/go-htmx/api"
	"github.com/ctheil/go-htmx/handlers"
)

func makeTodo(t api.Todo) string {
	var checked string
	if t.Complete {
		checked = "checked"
	} else {
		checked = ""
	}

	tStr := fmt.Sprintf(`
			<tr>
				<td>
					<input class="check" type="checkbox" hx-swap="none" hx-put="/complete?id=%v" %v>
				</td>
				<td>%v</td>
				<td>Today</td>
				<td>
					<button>Edit</button>
				</td>
			</tr>
			`, t.Id, checked, t.Name)
	return tStr
}
func getTodo(id uuid.UUID, todos []api.Todo) (t api.Todo, idx int) {
	idx = slices.IndexFunc(todos, func(t api.Todo) bool { return t.Id == id })
	return todos[idx], idx
}
func saveTodo(t api.Todo, idx int, tt []api.Todo) {
	tt[idx] = t
}

func main() {

	log.SetReportCaller(true)

	var r *chi.Mux = chi.NewRouter()
	handlers.Handler(r)

	fmt.Println("Starting GO API service")

	todos := []api.Todo{}

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		var out string
		for _, t := range todos {
			tStr := makeTodo(t)

			out += tStr
		}
		fmt.Println(out)
		fmt.Fprintf(w, out)
	})

	http.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Println("Error parsing form")
		}

		if r.Method == "POST" {
			name := r.FormValue("name")
			id, err := uuid.NewUUID()
			if err != nil {
				fmt.Println("UUID Error")

			}
			var t = api.Todo{Name: name, Id: id, Complete: false}
			todos = append(todos, t)

			tStr := makeTodo(t)
			fmt.Fprintf(w, tStr)
		}

	})

	http.HandleFunc("/complete", func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.URL.Query().Get("id"))
		if err != nil {
			fmt.Println("UUID parsing error")
		}
		t, idx := getTodo(id, todos)

		t.Complete = !t.Complete
		saveTodo(t, idx, todos)
		return
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Println("No data to parse")
		}

		q := r.FormValue("q")
		fmt.Fprintf(w, "<p>Result: %v</p>", q)

	})
	http.Handle("/", http.FileServer(http.Dir("./views")))

	if err := http.ListenAndServe(":42069", r); err != nil {
		fmt.Println(err)

	}

}
