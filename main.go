/*
We want to:
- Make requests for movies / tv for plex account
- have ability to see all requests through webpage
- "cross off" finished uploads
- request via webpage or rest call

Request:

{
    "title": string,
    "requestType": movie/tv show,
    "season" (optional): season #
}

*/

package main

import (
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
    "net/http"
    "os"
)

func AddRequest(w http.ResponseWriter, r *http.Request) {

}

type PlexMovieRequest struct {
    title       string
    requestType string
    season      string
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    mux := mux.NewRouter()
    mux.HandleFunc("/add", AddRequest).Methods("POST")
    // mux.HandleFunc("/subtract")

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Run(":" + port)

}
