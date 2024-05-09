FROM golang:1.22.2 as builder
 
# Create directory for alexandria app
WORKDIR /app

# Get dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy over alexandria files
COPY . ./

# Build binary
# RUN go build -o /usr/bin/alexandria-backend -v ./

# Expose port
EXPOSE 8080

# Start server on run
# ENTRYPOINT /usr/bin/alexandria-backend
