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
    "bytes"
    "github.com/golang/protobuf/jsonpb"
    "github.com/golang/protobuf/ptypes"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

const FileLocation = "movies.json"

var FilePerm os.FileMode = 0777

var unmarshaler = &jsonpb.Unmarshaler{}
var marshaler = &jsonpb.Marshaler{}

func AddRequest(w http.ResponseWriter, r *http.Request) {
    // todo: validate against the movie database, tmdb when you get a api key
    rq := PlexMovieRequest{}
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
    // rq.Uuid = uuid.New().String()
    rq.TimeRequested = ptypes.TimestampNow()

    // read file, then write to file
    // raw, err := ioutil.Reader(FILELOCATION)
    fileData, err := ioutil.ReadFile(FileLocation)
    if err != nil {
        log.Fatal("HELP ME", err)
    }
    raw := bytes.NewReader(fileData)
    requestList := RequestList{}
    if err = unmarshaler.Unmarshal(raw, &requestList); err != nil {
        log.Fatal("HELP ME w/ unmarshaler! ", err)
    }
    requestList.Shitwewant[uuid.New().String()] = &rq
    if err = WriteToFile(&requestList, FileLocation); err != nil {
        log.Fatal("BORKEN! ", err)
    }
    w.Write([]byte("hi"))
}

func WriteToFile(msg *RequestList, fileLoc string) error {
    file, err := os.Create(fileLoc)
    if err != nil {
        return err
    }
    if err = marshaler.Marshal(file, msg); err != nil {
        return err
    }
    return nil
}

func Validate(movieRequest *PlexMovieRequest) error {
    // fmt.Println(movieRequest)
    // validate w/ movie database
    return nil
}

func FinishHim(w http.ResponseWriter, r *http.Request) {
    SubtractRequest(w, r)
}

func SubtractRequest(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    itemId := params["id"]

    fileData, err := ioutil.ReadFile(FileLocation) // just pass the file name

    if err != nil {
        log.Println("couldn't file file at "+FileLocation, err)
    }

    byteReader := bytes.NewReader(fileData)
    currentReqs := &RequestList{}

    if err := unmarshaler.Unmarshal(byteReader, currentReqs); err != nil {
        log.Println("couldn't parse file", err)
    }

    _, ok := currentReqs.Shitwewant[itemId]

    if ok {
        delete(currentReqs.Shitwewant, itemId)
    } else {
        //it's not there it's not there
    }
    //overwrite file
    if err = WriteToFile(currentReqs, FileLocation); err != nil {
        log.Fatal("BORKEN! ", err)
    }
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
    mux.HandleFunc("/finishhim/{id}", FinishHim).Methods("POST")

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Run(":" + port)

}
