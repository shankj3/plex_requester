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
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
)

const FileLocation = "movies.json"
const FinishedLocation = "finished.json"

var FilePerm os.FileMode = 0777

var unmarshaler = &jsonpb.Unmarshaler{}
var marshaler = &jsonpb.Marshaler{}

func AddRequest(w http.ResponseWriter, r *http.Request) {
    // todo: validate against the movie database, tmdb when you get a api key
    tvSeason, err := strconv.ParseInt(r.FormValue("season"), 10, 32)

    if err != nil {
        log.Println("tv season wasn't a number", err)
    }

    rq := &PlexMovieRequest{
        Title:       r.FormValue("title"),
        RequestType: r.FormValue("requesttype"),
        Season:      int32(tvSeason),
    }

    // validate
    if err := Validate(rq); err != nil {
        log.Println("Missed Validations")
        //  return missing fields
    }
    // once validated, add timestamp and uuid
    // rq.Uuid = uuid.New().String()
    rq.TimeRequested = ptypes.TimestampNow()
    Append(rq, FileLocation, uuid.New().String())
    http.Redirect(w, r, "/", http.StatusSeeOther)

}

func Append(request *PlexMovieRequest, fileLoc string, uuid string) {
    requestList := &RequestList{}
    ReadFromFile(requestList, fileLoc)
    requestList.Shitwewant[uuid] = request
    if err := WriteToFile(requestList, fileLoc); err != nil {
        log.Fatal("error writing to file ", err)
    }
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

func ReadFromFile(requestList *RequestList, fileLoc string) {
    fileData, err := ioutil.ReadFile(fileLoc) // just pass the file name

    if err != nil {
        log.Println("couldn't file file at "+fileLoc, err)
    }

    byteReader := bytes.NewReader(fileData)

    if err := unmarshaler.Unmarshal(byteReader, requestList); err != nil {
        log.Println("couldn't parse file", err)
    }

}

func Validate(movieRequest *PlexMovieRequest) error {
    // fmt.Println(movieRequest)
    // validate w/ movie database
    return nil
}

func FinishHim(w http.ResponseWriter, r *http.Request) {
    r.Header.Set("FINISHED", "true")
    SubtractRequest(w, r)
}

func SubtractRequest(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    itemId := params["id"]

    currentReqs := &RequestList{}
    ReadFromFile(currentReqs, FileLocation)

    plexReq, ok := currentReqs.Shitwewant[itemId]

    if ok {
        delete(currentReqs.Shitwewant, itemId)
    } else {
        //it's not there it's not there
    }
    // if finished header set, add to finished.json before deleting from FileLocation
    log.Printf("REQUEST HEADERS!! %v", r.Header)
    if finished := r.Header.Get("FINISHED"); finished == "true" {
        Append(plexReq, FinishedLocation, itemId)
    }
    //overwrite file
    if err := WriteToFile(currentReqs, FileLocation); err != nil {
        log.Fatal("BORKEN! ", err)
    }
}

func ListRequests(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    file, err := ioutil.ReadFile(FileLocation)
    if err != nil {
        w.Write([]byte("cannot fulfill ListRequests"))
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write(file)
    }

}

func Homepage(w http.ResponseWriter, r *http.Request) {
    currentReqs := &RequestList{}
    ReadFromFile(currentReqs, FileLocation)
    renderTemplate(w, "index", currentReqs)
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    mux := mux.NewRouter()

    mux.HandleFunc("/", Homepage).Methods("GET")

    //TODO: validate adding input with tmdb api
    mux.HandleFunc("/add", AddRequest).Methods("POST")
    mux.HandleFunc("/subtract/{id}", SubtractRequest).Methods("DELETE")

    //TODO: email or text notification that something new got added?
    mux.HandleFunc("/finishhim/{id}", FinishHim).Methods("POST")

    mux.HandleFunc("/requests", ListRequests).Methods("GET")

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Run(":" + port)

}

//copied verbatim from golang docs
func renderTemplate(w http.ResponseWriter, tmpl string, currentList *RequestList) {
    t, _ := template.ParseFiles(tmpl + ".html")
    t.Execute(w, currentList)
}
