FROM golang:1.22-alpine as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o keepass-vault-sync

FROM scratch
COPY --from=build /app/keepass-vault-sync /keepass-vault-sync
ENTRYPOINT ["/keepass-vault-sync"]
