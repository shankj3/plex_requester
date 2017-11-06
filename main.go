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

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    mux := mux.NewRouter()

    //TODO: validate adding input with tmdb api
    mux.HandleFunc("/add", AddRequest).Methods("POST")

    // mux.HandleFunc("/subtract")

    //TODO: email or text notification that something new got added?
    // mux.HandleFunc("/finishhim")

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Run(":" + port)

}
