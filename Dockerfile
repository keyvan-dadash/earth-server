<<<<<<< HEAD
FROM golang AS builder
=======
FROM golang
>>>>>>> 21849edc12af8fa47588649c35ca44fb2ba458c8

WORKDIR /go/src/earth

COPY go.mod go.sum ./

RUN go mod download

COPY . .

<<<<<<< HEAD
RUN go build -o earth ./cmd/earth/

FROM ubuntu:latest  

WORKDIR /root/

COPY --from=builder /go/src/earth .

EXPOSE 8080

CMD ["./earth"]  
=======
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
>>>>>>> 21849edc12af8fa47588649c35ca44fb2ba458c8



