FROM golang:1.22.2 as builder

# Create directory for alexandria app
WORKDIR /app

# Copy over alexandria files
COPY . ./

# Build binary
# RUN go build -o /usr/bin/alexandria-backend -v ./

# Get missing dependencies
RUN go mod download
RUN go get github.com/golangci/golangci-lint

# Expose port
EXPOSE 8080

# Start server on run
# ENTRYPOINT /usr/bin/alexandria-backend
