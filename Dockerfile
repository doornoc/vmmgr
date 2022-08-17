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
RUN usermod -u 1000 root && groupmod -g 1000 root
COPY ./config /root/.ssh/config
COPY ./id_rsa /root/.ssh/id_rsa
COPY ./id_rsa.pub /root/.ssh/id_rsa.pub

CMD ["/backend", "start", "--config", "config.json"]