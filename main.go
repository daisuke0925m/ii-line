package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ii-line/handlers"
)

func main() {
	// ハンドラの登録
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/callback", handlers.LineHandler)

	// HTTPサーバを起動
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
