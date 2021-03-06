
FROM golang:1.14 as builder

ARG VERSION=
ARG GOARCH=
ARG GOOS=

WORKDIR /go/src/github.com/bartlettc22/wx200/
COPY ./ .

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -mod vendor -ldflags "-s -w -X main.version=${VERSION}" -v -a -o bin/wx200 .

FROM alpine:3.12

COPY --from=builder /go/src/github.com/bartlettc22/wx200/bin/wx200 /usr/bin

ENTRYPOINT ["wx200"]