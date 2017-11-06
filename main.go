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
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const FileLocation = "movies.json"

var unmarshaler = &jsonpb.Unmarshaler{}

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
	rq.Uuid = uuid.New().String()
	rq.TimeRequested = ptypes.TimestampNow()
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
	currentReqs := *RequestList{}

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
