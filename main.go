/*
We want to:
- Make requests for movies / tv for plex account
- have ability to see all requests through webpage
- "cross off" finished uploads
- request via webpage or rest call
*/

package main

import (
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
    // "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    mux := mux.NewRouter()
    // mux.HandleFunc("/add")
    // mux.HandleFunc("/subtract")

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Run(":" + port)

}
