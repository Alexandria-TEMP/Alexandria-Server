FROM golang:1.22.2 as builder

# Create directory for alexandria app
WORKDIR /app

# Copy over alexandria files
COPY . ./

# Get module dependencies
RUN go mod download
RUN go get github.com/golangci/golangci-lint@v1.58.1
RUN go install go.uber.org/mock/mockgen@v0.4.0

# Build binary
RUN go build -o /usr/bin/alexandria-backend -v ./

# Expose port
EXPOSE 8080

# Start server on run
ENTRYPOINT /usr/bin/alexandria-backend
