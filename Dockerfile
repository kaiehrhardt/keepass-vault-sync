FROM golang:1.22-alpine as build
WORKDIR /app
COPY . .
RUN go get -d -v ./... && \
    go install -v ./...

FROM alpine:latest as runtime
COPY --from=build /go/bin/keepass-vault-sync /usr/local/bin/
ENTRYPOINT ["keepass-vault-sync"]
