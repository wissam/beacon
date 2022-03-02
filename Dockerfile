# syntax=docker/dockerfile:1

FROM golang:1.17.7 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./  ./

RUN go build -o /beacon cmd/hookserver.go 

FROM golang:1.17.7 AS deploy


COPY --from=build /beacon  /beacon

EXPOSE 8080

ENTRYPOINT ["/beacon"]





