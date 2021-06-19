# Build stage
FROM golang:1.15-alpine3.13 AS build-env
WORKDIR /go/src/github.com/itiky/eth-block-proxy
RUN apk add --no-cache make bash git build-base
COPY . .

RUN make install

# Run stage
FROM alpine:3.13

EXPOSE 2412

RUN apk --no-cache add ca-certificates
WORKDIR /opt/app
COPY --from=build-env /go/bin/eth-block-proxy .

CMD ["/opt/app/eth-block-proxy", "server", "--log-level=debug"]
