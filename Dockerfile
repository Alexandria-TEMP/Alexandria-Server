FROM golang:1.22.2 as build
 
ARG GOPKG
ARG GOBIN

# Create directory for alexandria app
WORKDIR /app

# Copy over alexandria files
COPY . ./

# Get missing dependencies
RUN go mod download

# Developer tools
# TODO for prod these can be removed
RUN go get github.com/golangci/golangci-lint
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Build binary
RUN go build -o /usr/bin/alexandria-backend -v ./

# Expose port
EXPOSE 8080

FROM build AS run

# Start server on run
ENTRYPOINT /usr/bin/alexandria-backend
