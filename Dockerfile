#############################################################

FROM --platform=linux/x86_64 golang:1.19.0-alpine3.16 as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOARCH=amd64 \
    GOOS=linux
WORKDIR /pkg
COPY go.mod go.sum ./
RUN go mod download
COPY api ./api
COPY cmd ./cmd
COPY docs ./docs
COPY internal ./internal
RUN go build -o server ./cmd/app/main.go

#############################################################

FROM --platform=linux/x86_64 alpine:3.16.2 AS app
ENV GIN_MODE=release
WORKDIR /bin
COPY --from=builder /pkg/server /bin/server
CMD /bin/server