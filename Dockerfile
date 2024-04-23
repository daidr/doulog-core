############################
# STEP 1 build executable binary
############################
FROM golang:1.21-alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache 'git=~2'

# Install dependencies
ENV GO111MODULE=on
WORKDIR $GOPATH/src/packages/doulog-core/
COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/main .

############################
# STEP 2 build a small image
############################
FROM alpine:3

WORKDIR /

# Copy our static executable.
COPY --from=builder /go/main /go/main


ENV GIN_MODE release
ENV DOULOG_SERVER_PORT 8080
ENV DOULOG_SERVER_SITE http://localhost:3000
ENV DOULOG_PGSQL_HOST localhost
ENV DOULOG_PGSQL_PORT 5432
ENV DOULOG_PGSQL_USERNAME doulog
ENV DOULOG_PGSQL_PASSWORD doulog
ENV DOULOG_PGSQL_DATABASE doulog
ENV DOULOG_REDIS_HOST localhost
ENV DOULOG_REDIS_PORT 6379
ENV DOULOG_REDIS_AUTH ""
ENV DOULOG_AUTH_FRONTEND_CALLBACK_PREFIX "http://localhost:5432"

ENV DOULOG_AUTH_GITHUB_CLIENT_ID ""
ENV DOULOG_AUTH_GITHUB_CLIENT_SECRET ""

ENV DOULOG_LIMIT_MEDIA_FIEL_SIZE 51200
ENV DOULOG_LIMIT_MEDIA_MEDIA_SIZE 7680

EXPOSE 8080

WORKDIR /Go

RUN mkdir /go/data

VOLUME ["/go/data"]

# Run the Go Gin binary.
ENTRYPOINT ["/go/main"]