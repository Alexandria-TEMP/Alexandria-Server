# FROM golang:1.22.2 as build
FROM janneskelso/go-with-quarto:alpha as build
 
ARG GOPKG
ARG GOBIN

# Create directory for alexandria app
WORKDIR /app

# Copy over alexandria files
COPY . ./

# Get module dependencies
RUN go mod download
RUN go mod download

# Developer tools
# TODO for prod these can be removed
RUN go get github.com/golangci/golangci-lint
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN go install go.uber.org/mock/mockgen@v0.4.0

# Build binary
RUN go build -o /usr/bin/alexandria-backend -v ./

# Generate API spec
RUN swag init -g alexandria.go

# Expose port
EXPOSE 8080

FROM build AS run

# Start server on run
ENTRYPOINT /usr/bin/alexandria-backend
