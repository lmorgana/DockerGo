package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
	"os"
)

type Entry struct {
	Number int
	Double int
	Square int
}

var DATA []Entry
var tFile string

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Host: %s Path: %s\n", r.Host, r.URL.Path)
	myT := template.Must(template.ParseGlob(tFile))
	myT.ExecuteTemplate(w, tFile, DATA)
}

func main() {

	arguments := os.Args
	if len(arguments) != 4 {
		fmt.Println("Need DB: username and password + Template File!")
		return
	}
	connStr := "user=postgres password=postgres host=host.docker.internal port=5432 dbname=postgres sslmode=disable" //postgresql://postgres:postgres@localhost:5432/db sslmode=disable"
	tFile = arguments[3]

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	fmt.Println("make db")
	_, err = db.Exec("CREATE TABLE data (number INTEGER PRIMARY KEY, double INTEGER, square INTEGER )")
	if err != nil && err.Error() != "pq: relation \"data\" already exists" {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Emptying database table.")
	_, err = db.Exec("DELETE FROM data")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Populating database")
	sqlStatement := `INSERT INTO data (number, double, square)
VALUES ($1, $2, $3)`

	for i := 20; i < 50; i++ {
		_, _ = db.Exec(sqlStatement, i, 2*i, i*i)
	}

	rows, err := db.Query("SELECT * FROM data")
	if err != nil {
		fmt.Println(nil)
		return
	}
	var n int
	var d int
	var s int
	for rows.Next() {
		err = rows.Scan(&n, &d, &s)
		temp := Entry{Number: n, Double: d, Square: s}
		DATA = append(DATA, temp)
	}

	http.HandleFunc("/", myHandler)
	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
