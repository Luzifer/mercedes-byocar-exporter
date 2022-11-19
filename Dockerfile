FROM golang:alpine as builder

COPY . /src/mercedes-byocar-exporter
WORKDIR /src/mercedes-byocar-exporter

RUN set -ex \
 && apk add --update \
      build-base \
      git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath

FROM alpine:latest

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/mercedes-byocar-exporter /usr/local/bin/mercedes-byocar-exporter

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/mercedes-byocar-exporter"]
CMD ["--"]

# vim: set ft=Dockerfile:
