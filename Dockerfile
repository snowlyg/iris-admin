# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.

FROM golang:1.14.1

LABEL maintainer="snowlyg <569616226@qq.com>"

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

# Copy the local package files to the container's workspace.
COPY ./config /go/src/github.com/snowlyg/IrisAdminApi/config
COPY ./controllers /go/src/github.com/snowlyg/IrisAdminApi/controllers
COPY ./libs /go/src/github.com/snowlyg/IrisAdminApi/libs
COPY ./middleware  /go/src/github.com/snowlyg/IrisAdminApi/middleware
COPY ./models  /go/src/github.com/snowlyg/IrisAdminApi/models
COPY ./routes  /go/src/github.com/snowlyg/IrisAdminApi/routes
COPY ./seeder  /go/src/github.com/snowlyg/IrisAdminApi/seeder
COPY ./sysinit  /go/src/github.com/snowlyg/IrisAdminApi/sysinit
COPY ./transformer  /go/src/github.com/snowlyg/IrisAdminApi/transformer
COPY ./validates  /go/src/github.com/snowlyg/IrisAdminApi/validates
COPY ./web_server  /go/src/github.com/snowlyg/IrisAdminApi/web_server
COPY ./main.go  /go/src/github.com/snowlyg/IrisAdminApi/main.go
COPY ./www/dist  /go/src/github.com/snowlyg/IrisAdminApi/www/dist
COPY ./go.mod  /go/src/github.com/snowlyg/IrisAdminApi/go.mod
COPY ./go.sum  /go/src/github.com/snowlyg/IrisAdminApi/go.sum

#build the application
RUN cd /go/src/github.com/snowlyg/IrisAdminApi && \
     go get -u github.com/go-bindata/go-bindata/v3/go-bindata && \
     go generate && \
     go build -o main

# Run the command by default when the container starts.
ENTRYPOINT /go/src/github.com/snowlyg/IrisAdminApi/main

# Document that the service listens on port 8085
EXPOSE 8085
