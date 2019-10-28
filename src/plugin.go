package main

import (
        "net/http"
)

func main() {
        // Each function and its helpers are defined in a separate .go file
        http.HandleFunc("/evaluate", Evaluate)
        http.HandleFunc("/status", Status)
        http.HandleFunc("/run", Run)

        http.ListenAndServe(":3000", nil)
}

