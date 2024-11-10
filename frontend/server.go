package main

import (
    "log"
    "net/http"
)

func main() {
    fs := http.FileServer(http.Dir("./frontend"))
    http.Handle("/", fs)

    log.Println("Frontend server running on http://localhost:3000")
    err := http.ListenAndServe(":3000", nil)
    if err != nil {
        log.Fatal(err)
    }
}
