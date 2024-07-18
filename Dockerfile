FROM docker.io/janneskelso/go-with-quarto:alpha as build
 
ENV GIN_MODE=release

ARG GOPKG
ARG GOBIN

# Create directory for alexandria app
WORKDIR /app

# Set git config
RUN git config --global user.name "Alexandria Bot"
# TODO change email
RUN git config --global user.email "todo@todo.todo" 

# Copy over alexandria files
COPY . ./

# Get module dependencies
RUN go mod download

# Install developer tools
# TODO for prod these can be removed
RUN go get github.com/golangci/golangci-lint
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN go install go.uber.org/mock/mockgen@v0.4.0

# Generate API spec
RUN swag init -g alexandria.go 

# Build binary
RUN go build -o /usr/bin/alexandria-backend -v ./

# Expose port
EXPOSE 8080

FROM build AS run

# Start server on run
ENTRYPOINT /usr/bin/alexandria-backend
