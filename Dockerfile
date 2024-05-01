FROM golang:1.22.2 as builder
 
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o ./bin/alexandria-backend -v ./

EXPOSE 8080

ENTRYPOINT /app/bin/alexandria-backend