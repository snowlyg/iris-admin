# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.

FROM archlinux

LABEL maintainer="snowlyg <569616226@qq.com>"

COPY ./main_lin /go/src/github.com/snowlyg/blog/main_lin
COPY ./application.example.yml /go/src/github.com/snowlyg/blog/application.yml
COPY ./rbac_model.conf /go/src/github.com/snowlyg/blog/rbac_model.conf
COPY ./seeder/data /go/src/github.com/snowlyg/blog/seeder/data


# Run the command by default when the container starts.
ENTRYPOINT /go/src/github.com/snowlyg/blog/main_lin

# Document that the service listens on port 80
EXPOSE 80
