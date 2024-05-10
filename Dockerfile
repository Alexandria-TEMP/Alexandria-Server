FROM golang:1.22.2 as builder
 
ARG GOPKG
ARG GOBIN

# Create directory for alexandria app
WORKDIR /app

# Copy over alexandria files
COPY . ./

# Copy ssh key
COPY ./.devcontainer/id_rs[a] /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa

# Build binary
# RUN go build -o /usr/bin/alexandria-backend -v ./

# Get missing dependencies
RUN go mod download
RUN go install github.com/golangci/golangci-lint

# Expose port
EXPOSE 8080

# Start server on run
# ENTRYPOINT /usr/bin/alexandria-backend
