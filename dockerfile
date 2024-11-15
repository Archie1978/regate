FROM golang:1.23 AS builder

WORKDIR /src



#COPY go.* /src
#COPY grdp/go.* /src/grdp
#RUN /bin/sh -c "git config --global --add safe.directory /src && cd /src && go mod download"

COPY . .
RUN /bin/sh -c  "git config --global --add safe.directory /src &&cd /src/cmd/regate-daemon && go build"

FROM ubuntu:24.04

WORKDIR /app

COPY cmd/regate-daemon/configuration.json.example /app/configuration.json
COPY --from=builder --chmod=0755 /src/cmd/regate-daemon/regate-daemon /app/regate-daemon

EXPOSE 4203/tcp
CMD /app/regate-daemon
