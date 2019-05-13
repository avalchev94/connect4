package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/avalchev94/tarantula"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	if err := http.ListenAndServe(*addr, tarantula.NewServer()); err != nil {
		log.Fatalln(err)
	}
}

var db *sql.DB

func InitDB(driver, dataSourceName string) error {
	var err error
	db, err = sql.Open(driver, dataSourceName)
	return err
}
