package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("This is Spoker!")

    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Setting up server")
    })

    http.ListenAndServe(":8080", nil)
}
