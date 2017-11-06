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
    "io/ioutil"
    "fmt"
    "github.com/golang/protobuf/jsonpb"
    "bytes"
)

const FileLocation = "file.json"
var unmarshaler = &jsonpb.Unmarshaler{}

func AddRequest(w http.ResponseWriter, r *http.Request) {

}


func SubtractRequest(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    itemId := params["id"]

    fileData, err := ioutil.ReadFile(FileLocation) // just pass the file name

    if err != nil {
        fmt.Print(err)
    }

    byteReader := bytes.NewReader(fileData)
    currentReqs := &RequestList{}

    if err := unmarshaler.Unmarshal(byteReader, currentReqs); err != nil {
        fmt.Print(err)
    }

    //for index, item := range people {
    //if item.ID == params["id"] {
    //people = append(people[:index], people[index+1:]...)
    //break
    //}
    //json.NewEncoder(w).Encode(people)
    //}

}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    mux := mux.NewRouter()

    //TODO: validate adding input with tmdb api
    mux.HandleFunc("/add", AddRequest).Methods("POST")

    mux.HandleFunc("/subtract/{id}", SubtractRequest).Methods("DELETE")

    //TODO: email or text notification that something new got added?
    // mux.HandleFunc("/finishhim")

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Run(":" + port)

}
