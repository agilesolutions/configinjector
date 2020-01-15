# Multistage GO build and package
FROM golang:alpine AS builder

ADD . /src

WORKDIR /src

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o configinjector .

# Final stage packaging the executable
FROM alpine:latest

LABEL maintainer="Robert Rong <robert.rong@agile-solutions.ch>"

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder /src/configinjector /app/

ENTRYPOINT ["./configinjector"]
