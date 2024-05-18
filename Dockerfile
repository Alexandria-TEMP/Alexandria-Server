FROM golang:1.22.2 as build
 
ARG GOPKG
ARG GOBIN

# Create directory for alexandria app
WORKDIR /app

# Copy over alexandria files
COPY . ./

# Get missing dependencies
RUN go mod download
RUN go get github.com/golangci/golangci-lint

# Build binary
RUN go build -o /usr/bin/alexandria-backend -v ./

# Expose port
EXPOSE 8080

FROM build AS run

# Start server on run
ENTRYPOINT /usr/bin/alexandria-backend
