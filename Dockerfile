FROM golang:1.17

WORKDIR /app
ADD . /app 

RUN go get -d -v ./...
RUN go install -v ./...

CMD go run . 