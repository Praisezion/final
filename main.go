package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"final/pkg/db"
)

func main() {
	// server prot
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// BD  proverka
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()

	fmt.Printf("База данных: %s\n", dbFile)

	webDir := "./web"

	// proverka dir
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatal("Директория web не найдена")
	}

	fs := http.FileServer(http.Dir(webDir))

	http.Handle("/", fs)

	// start servera
	addr := ":" + port
	fmt.Printf("Сервер запущен на http://localhost:%s\n", port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
