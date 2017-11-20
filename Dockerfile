FROM golang:1.9.2


WORKDIR /go/src/plex_requester
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."
ENTRYPOINT ["/go/bin/plex_requester"]
