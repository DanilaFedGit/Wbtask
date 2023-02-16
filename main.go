package main

import (
	"fmt"
	"log"
	"net/http"
	db "project/db"
	jsoncheck "project/jscheck"
	"text/template"

	stan "github.com/nats-io/stan.go"
)

func main() {
	db_connect, err := db.ConnectDB(db.Config{"localhost", "5432", "qwerty", "qwerty", "postgres", "disable"})
	if err != nil {
		log.Fatal(err)
	}
	cache, err := db.GetData(db_connect)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", main_page)
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {

		check_id := r.FormValue("id")
		if _, ok := cache[check_id]; !ok {
			fmt.Fprintf(w, "Не существуют данные с id: %s", check_id)
		} else {
			tmpl, _ := template.ParseFiles("template/check.html")
			tmpl.Execute(w, jsoncheck.Getjson(cache[check_id]))

		}
	})
	sc, err := stan.Connect("test-cluster", "client1")
	if err != nil {
		log.Fatal(err)
	}
	sub, err := sc.Subscribe("foo1", func(m *stan.Msg) {
		id, data, err := jsoncheck.ValidCheck(string(m.Data))
		if err == nil {
			if _, ok := cache[id]; !ok {
				db.InsertData(db_connect, id, data)
				cache[id] = data
			}
		}
	}, stan.DeliverAllAvailable())
	defer sub.Unsubscribe()
	if err != nil {
		log.Fatal(err)
	}
	http.ListenAndServe(":8080", nil)
}
func main_page(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/index.html")
}
