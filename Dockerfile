FROM golang:1.14 AS base
RUN set -ex \
 && go get github.com/mitchellh/gox \
 && go get github.com/vektra/mockery/.../

FROM base AS builder
COPY . /go/src/github.com/xsteadfastx/jitsiexporter
WORKDIR /go/src/github.com/xsteadfastx/jitsiexporter
RUN set -ex \
 && make build

FROM scratch
COPY --from=builder /go/src/github.com/xsteadfastx/jitsiexporter/jitsiexporter_linux_amd64 /bin/jitsiexporter
EXPOSE 6700
ENTRYPOINT ["/bin/jitsiexporter", "-debug=true", "-host=0.0.0.0"]
