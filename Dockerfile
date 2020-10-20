# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.

FROM archlinux

LABEL maintainer="snowlyg <569616226@qq.com>"

COPY ./main_lin  /go/src/github.com/snowlyg/IrisAdminApi/main_lin
COPY ./seeder/data  /go/src/github.com/snowlyg/IrisAdminApi/seeder/data
COPY ./application.yml  /go/src/github.com/snowlyg/IrisAdminApi/application.yml
COPY ./rbac_model.conf  /go/src/github.com/snowlyg/IrisAdminApi/rbac_model.conf

# Run the command by default when the container starts.
ENTRYPOINT /go/src/github.com/snowlyg/IrisAdminApi/main_lin

# Document that the service listens on port 8085
EXPOSE 80
