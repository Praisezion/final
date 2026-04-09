package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"final/pkg/api"
	"final/pkg/db"
)

func main() {
	// server port ***
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// DB check ***
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()

	fmt.Printf("База данных: %s\n", dbFile)

	// api
	api.Init()

	webDir := "./web"
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatal("Директория web не найдена")
	}

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	// start server
	addr := ":" + port
	fmt.Printf("start http://localhost:%s\n", port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("errpr:", err)
	}
}
