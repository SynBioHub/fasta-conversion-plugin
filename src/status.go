package main

import (
        "fmt"
        "net/http"
)

func Status(w http.ResponseWriter, r *http.Request) {
        // We are ready to go as long as we can connect to the validator
        _, err := http.Get("https://validator.sbolstandard.org/")
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        } else {
                fmt.Fprintln(w, "Ready to go!")
        }
}
