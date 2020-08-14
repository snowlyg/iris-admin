# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.

FROM golang:1.14.1
# Copy the local package files to the container's workspace.
ADD server /go/src/github.com/snowlyg/IrisAdminApi/

#build the application
RUN cd /go/src/github.com/snowlyg/IrisAdminApi/ go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go build

# Run the command by default when the container starts.
ENTRYPOINT /go/src/github.com/snowlyg/IrisAdminApi/server

# Document that the service listens on port 8081
EXPOSE 8085
