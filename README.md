# Alexandria Backend



## Getting started

### Installing Go
Install go version 1.22.2 from https://go.dev/dl/ and then set your environment variables as follows:
 - Add the go binary directory to PATH 
 - Add the go installation directory to GOROOT
 - Add the alexandria-backend parent directory (ie the directory above this README) as GOPATH

### Setting up the Linter in VSCode
In order to lint locally install golangci-lint by running the following command in the root directory:
```go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2```

Then in the VSCode workspace settings:
 - set *Lint Tool* to *golangci-lint*
 - set *Lint on Save* to *workspace*
 - add ```--config=${workspaceFolder}/.golangci.yml``` as a *Lint Flag*

### Go CLI Commands
 - To get all dependencies run ```go mod tidy``` in this directory.    
 - To run main run ```go run ./...``` or ```go run ./cmd/alexandria/alexandria.go```     
 - To run all tests run ```go test ./...``` or ```go test ./test```
 - To build the application run ```go build ./cmd/alexandria/alexandria.go``` or ```go build ./...```

### Deployment
SSH into ```ssh [NetID]@student-linux.tudelft.nl``` and when prompted enter your SSO password.