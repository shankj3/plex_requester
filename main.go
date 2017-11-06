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
    "bytes"
    "fmt"
    "github.com/golang/protobuf/jsonpb"
    "github.com/golang/protobuf/ptypes"
    // "github.com/google/uuid"
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
    w.Write([]byte("hi"))
}

func Validate(movieRequest *PlexMovieRequest) error {
    // fmt.Println(movieRequest)
    // validate w/ movie database
    return nil
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

	_, ok := currentReqs.Shitwewant[itemId]

	if ok {
		delete(currentReqs.Shitwewant, itemId)
	} else {
		//it's not there it's not there
	}

	//overwrite file
	reqs, err := marshaler.MarshalToString(currentReqs)
	if err != nil {
		fmt.Print(err)
	}
	err = ioutil.WriteFile(FileLocation, []byte(reqs), FilePerm)
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
