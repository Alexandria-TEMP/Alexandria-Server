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
 - To get all dependencies run ```go mod download``` in this directory.    
 - To run main run ```go run``` or ```go run alexandria.go```     
 - To run all tests run ```go test ./...``` or ```go test ./test```
 - To build the application run ```go build``` or ```go build alexandria.go```

### Deployment
SSH into ```ssh [NetID]@student-linux.tudelft.nl``` and when prompted enter your SSO password.

### Docker
 - To build the go docker container run ```docker build -t alexandria-backend .```
 - To start the container run ```docker run alexandria-backend```
 - To check if the container is running with the server run ```docker ps```

## Default readme info to review and consider

Use the built-in continuous integration in GitLab.

- [ ] [Get started with GitLab CI/CD](https://docs.gitlab.com/ee/ci/quick_start/index.html)
- [ ] [Analyze your code for known vulnerabilities with Static Application Security Testing (SAST)](https://docs.gitlab.com/ee/user/application_security/sast/)
- [ ] [Deploy to Kubernetes, Amazon EC2, or Amazon ECS using Auto Deploy](https://docs.gitlab.com/ee/topics/autodevops/requirements.html)
- [ ] [Use pull-based deployments for improved Kubernetes management](https://docs.gitlab.com/ee/user/clusters/agent/)
- [ ] [Set up protected environments](https://docs.gitlab.com/ee/ci/environments/protected_environments.html)
