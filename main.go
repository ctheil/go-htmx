package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

func formatSuffix(out *string, num int) {
	*out = "th"
	if num == 1 {
		*out = "st"
	} else if num%2 == 0 {
		*out = "nd"
	} else if num%3 == 0 {
		*out = "rd"
	}
}

type Todo struct {
	Name     string
	Id       uuid.UUID
	Complete bool
}

func main() {
	todos := []Todo{}

	//r.HandleFunc("/", controller.Home).Methods("GET")
	http.HandleFunc("/swap", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Made it to /home")

		//fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
		if r.Method == "GET" {

			fmt.Fprintf(w, "<p>GET</p>")
		}
		if r.Method == "POST" {

			fmt.Fprintf(w, "<p>POST</p>")
		}

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

			var numSuffix string
			formatSuffix(&numSuffix, len(todos))

			fmt.Println(t)
			var checked string
			if t.Complete {
				checked = "checked"
			} else {
				checked = ""
			}
			fmt.Fprintf(w, `
			<li>
			<div class="todo-item">

			<input class="check" type="checkbox" role="switch" hx-post="/complete?id=%v" %v>
			%v%v Todo: %v
			</div>
			</li>
			`, t.Id, checked, len(todos), numSuffix, t.Name)

		}

	})

	http.HandleFunc("/complete", func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.URL.Query().Get("id"))
		if err != nil {
			log.Fatal("uuid parse error")
		}
		fmt.Println(id)

		idx := slices.IndexFunc(todos, func(t Todo) bool { return t.Id == id })

		t := todos[idx]
		t.Complete = true

		fmt.Println(t)
		fmt.Fprintf(w, `
		<li>
			<div class="todo-item">

			<input class="check" type="checkbox" role="switch" hx-post="/complete?id=%v" %v>
			%v%v Todo: %v
			</div>
		</li>
		`, id, "checked", idx, "", t.Name)

	})
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Fatal("No data to parse")
		}

		q := r.FormValue("q")
		fmt.Println(q)
		fmt.Fprintf(w, "<p>Result: %v</p>", q)

	})
	http.Handle("/", http.FileServer(http.Dir("./views")))

	if err := http.ListenAndServe(":42069", nil); err != nil {
		fmt.Println(err)
		log.Fatal("Error")
	}

}
