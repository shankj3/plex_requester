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
    "season" (optional): season #,
}

*/

package main

import (
    // "fmt"
    "bufio"
    "github.com/golang/protobuf/jsonpb"
    "github.com/golang/protobuf/ptypes"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
    "log"
    "net/http"
    "os"
)

const FILELOCATION = "/Users/jesseshank/go/src/github.com/shankj3/plex_requester/requests/requestList.json"

func AddRequest(w http.ResponseWriter, r *http.Request) {
    // todo: validate against the movie database, tmdb when you get a api key
    rq := []PlexMovieRequest{}
    unmarshaler := &jsonpb.Unmarshaler{
        AllowUnknownFields: true,
    }
    if err := unmarshaler.Unmarshal(r.Body, &rq); err != nil {
        log.Println("Can't unmarshal object!", err)
        // return error summary and erro code
    }
    // validate
    if err := Validate(&rq); err != nil {
        log.Println("Missed Validations")
        //  return missing fields
    }
    // once validated, add timestamp and uuid
    rq.Uuid = uuid.New().String()
    rq.TimeRequested = ptypes.TimestampNow()

    // read file, then write to file
    // raw, err := ioutil.Reader(FILELOCATION)
    f, err := os.Open(FILELOCATION)
    if err != nil {
        log.Fatal(err)
    }
    raw := bufio.NewReader(f)
    requestList := RequestList{}
    if err = unmarshaler.Unmarshal(raw, &requestList); err != nil {
        log.Fatal(err)
    }
    log.Println(requestList)
    w.Write([]byte("hi"))

}

func Validate(movieRequest *PlexMovieRequest) error {
    // fmt.Println(movieRequest)
    // validate w/ movie database
    return nil
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
