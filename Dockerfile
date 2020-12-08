FROM golang AS builder

WORKDIR /go/src/earth

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o earth ./cmd/earth/

FROM ubuntu:latest  

WORKDIR /root/

COPY --from=builder /go/src/earth .

EXPOSE 8080

CMD ["./earth"]  



