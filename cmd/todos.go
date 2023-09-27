package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"github.com/ctheil/go-htmx/cmd/api"

)


func makeTodo(t Todo) string {
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
func getTodo(id uuid.UUID, todos []Todo) (t Todo, idx int){
		idx = slices.IndexFunc(todos, func(t Todo) bool { return t.Id == id })
		return todos[idx], idx
}
func saveTodo(t Todo, idx int, tt []Todo) {
	tt[idx] = t;
}


func main() {
	todos := []Todo{}

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		var out string
	for _, t := range todos {
		tStr := makeTodo(t);

		out += tStr;
	}
	fmt.Println(out)
	fmt.Fprintf(w, out)
	})

	http.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Fatal("No data to parse")
		}

		if r.Method == "POST" {
			name := r.FormValue("name")
			id, err := uuid.NewUUID()
			if err != nil {
				log.Fatal("UUID Error")

			}
			var t = Todo{Name: name, Id: id, Complete: false}
			todos = append(todos, t)

			tStr := makeTodo(t)
		fmt.Fprintf(w, tStr)
		}

	})

	http.HandleFunc("/complete", func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.URL.Query().Get("id"))
		if err != nil {
			log.Fatal("uuid parse error")
		}
		t, idx := getTodo(id, todos)

		t.Complete = !t.Complete
		saveTodo(t,idx, todos)
		return
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Fatal("No data to parse")
		}

		q := r.FormValue("q")
		fmt.Fprintf(w, "<p>Result: %v</p>", q)

	})
	http.Handle("/", http.FileServer(http.Dir("./views")))

	if err := http.ListenAndServe(":42069", nil); err != nil {
		fmt.Println(err)
		log.Fatal("Error")
	}

}
