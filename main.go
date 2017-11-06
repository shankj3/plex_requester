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
    //TODO: validate adding input with tmdb api
    // mux.HandleFunc("/add")

    //TODO: email or text notification that something new got added?
    // mux.HandleFunc("/subtract")


    // mux.HandleFunc("/finishhim")

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Run(":" + port)

}
