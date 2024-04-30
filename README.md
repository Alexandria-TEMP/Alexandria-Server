# Alexandria Backend



## Getting started

### Installing Go
Install go version 1.22.2 from https://go.dev/dl/
Ensure that the Go binary is in added to your PATH environemnt variable. Also set your GOPATH environment variable to the alexandria-backend directory (ie the directory of this README).

### Setting up the Linter in VSCode
In order to lint locally install golangci-lint by running the following command in the root directory:
```go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2```

Then in the VSCode workspace settings:
 - set *Lint Tool* to *golangci-lint*
 - set *Lint on Save* to *workspace*
 - add ```-c .golangci.yml``` as a *Lint Flag*
