FROM golang:1.18 AS builder
WORKDIR /build

COPY . .

RUN go build -o ./bin/sc-bot ./cmd/

FROM ubuntu:latest
WORKDIR /srv

COPY --from=builder /build/bin/sc-bot ./sc-bot
COPY ./database.json .

ENTRYPOINT [ "/srv/sc-bot" ]
