#
#  Compile interface ( nodejs vue )
#
#
FROM node:22-alpine AS builder-ui

WORKDIR /src
COPY webservice/www/regate/package.json .
RUN yarn install

COPY webservice/www/regate .
RUN yarn build


#
# Compile src go
#
FROM golang:1.23 AS builder

WORKDIR /src

COPY go.* /src
RUN mkdir /src/grdp
COPY grdp/go.* /src/grdp
RUN /bin/sh -c "git config --global --add safe.directory /src && cd /src && go mod download"

# get go src
COPY . .

# get UI
COPY --from=builder-ui /src/dist /src/webservice/www/regate/dist


RUN /bin/sh -c  "git config --global --add safe.directory /src &&cd /src/cmd/regate-daemon && go build"

FROM ubuntu:24.04

WORKDIR /app

COPY cmd/regate-daemon/configuration.json.example /app/configuration.json
COPY --from=builder --chmod=0755 /src/cmd/regate-daemon/regate-daemon /app/regate-daemon

EXPOSE 4203/tcp
CMD /app/regate-daemon
