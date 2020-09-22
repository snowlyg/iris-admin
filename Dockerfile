# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.

FROM golang:1.14

ENV  GO111MODULE=on
ENV  GOPROXY=https://goproxy.cn,direct

LABEL maintainer="snowlyg <569616226@qq.com>"

# Copy the local package files to the container's workspace.
COPY ./config /go/src/github.com/snowlyg/IrisAdminApi/config
COPY ./controllers /go/src/github.com/snowlyg/IrisAdminApi/controllers
COPY ./libs /go/src/github.com/snowlyg/IrisAdminApi/libs
COPY ./middleware /go/src/github.com/snowlyg/IrisAdminApi/middleware
COPY ./models /go/src/github.com/snowlyg/IrisAdminApi/models
COPY ./routes /go/src/github.com/snowlyg/IrisAdminApi/routes
COPY ./seeder /go/src/github.com/snowlyg/IrisAdminApi/seeder
COPY ./sysinit /go/src/github.com/snowlyg/IrisAdminApi/sysinit
COPY ./transformer /go/src/github.com/snowlyg/IrisAdminApi/transformer
COPY ./validates /go/src/github.com/snowlyg/IrisAdminApi/validates
COPY ./web_server /go/src/github.com/snowlyg/IrisAdminApi/web_server
COPY ./main.go /go/src/github.com/snowlyg/IrisAdminApi/main.go
COPY ./bindata.go /go/src/github.com/snowlyg/IrisAdminApi/bindata.go
COPY ./application.example.yml /go/src/github.com/snowlyg/IrisAdminApi/application.yml
COPY ./rbac_model.conf /go/src/github.com/snowlyg/IrisAdminApi/rbac_model.conf
COPY ./go.mod /go/src/github.com/snowlyg/IrisAdminApi/go.mod
COPY ./go.sum /go/src/github.com/snowlyg/IrisAdminApi/go.sum

#build the application
RUN cd /go/src/github.com/snowlyg/IrisAdminApi && \
     go build  -o main

# Run the command by default when the container starts.
ENTRYPOINT /go/src/github.com/snowlyg/IrisAdminApi/main

# Document that the service listens on port 8085
EXPOSE 80
