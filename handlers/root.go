package handlers

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	msg := "server is running"
	fmt.Fprintf(w, msg)
}
