ARG VERSION=1.24

FROM golang:${VERSION}-alpine AS build

SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

RUN apk add --update --no-cache \
    git \
    make \
    gcc alpine-sdk


WORKDIR $GOPATH/src/package

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY go.sum go.sum
COPY go.mod go.mod
COPY main.go main.go

# build binary
RUN go install

FROM alpine:3.21

# Containers must run as nonroot
ARG UID=10001
ARG USER=nonexistent

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/home/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

USER nonexistent

COPY --from=build /go/bin/kaet /go/bin/kaet
ENTRYPOINT [ "/go/bin/kaet" ]
