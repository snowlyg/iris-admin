# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
LABEL maintainer="snowlyg <569616226@qq.com>"

FROM golang:1.14.1
# Copy the local package files to the container's workspace.
ADD ./ /go/src/github.com/snowlyg/IrisAdminApi

#build the application
RUN cd /go/src/github.com/snowlyg/IrisAdminApi && \
     go env -w GO111MODULE=on && \
     go env -w GOPROXY=https://goproxy.cn,direct && \
     go get -u github.com/go-bindata/go-bindata/v3/go-bindata && \
     go generate && \
     go build -o main

# Run the command by default when the container starts.
ENTRYPOINT /go/src/github.com/snowlyg/IrisAdminApi/main

# Document that the service listens on port 8085
EXPOSE 8085
