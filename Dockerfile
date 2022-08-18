## Build
FROM golang:1.19-bullseye AS build

WORKDIR /app
COPY . ./
RUN go mod download
WORKDIR /app/cmd/backend
RUN apt-get update
RUN apt-get install -y libvirt-dev
RUN go build -o /backend


## Deploy
FROM ubuntu:22.04

WORKDIR /
COPY --from=build /backend /backend
RUN apt-get update
RUN apt-get install -y libvirt-dev ssh

CMD ["/backend", "start", "--config", "config.json"]
