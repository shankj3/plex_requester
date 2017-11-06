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
    "github.com/golang/protobuf/jsonpb"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
    "log"
    "net/http"
    "os"
)

func AddRequest(w http.ResponseWriter, r *http.Request) {
    // todo: validate against the movie database, tmdb when you get a api key
    rq := PlexMovieRequest{}
    unmarshaler := &jsonpb.Unmarshaler{
        AllowUnknownFields: true,
    }
    if err := unmarshaler.Unmarshal(r.Body, &rq); err != nil {
        log.Fatal("DISTRESS CALL!!!", err)
    }
    w.Write([]byte("hi"))

}

func Validate(title *string) error {
    return nil
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
