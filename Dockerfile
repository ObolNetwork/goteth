# syntax=docker/dockerfile:1
FROM golang:1.26.10-alpine as builder
RUN apk add --update git gcc g++ openssh-client make
WORKDIR /app
COPY go.mod go.sum ./
COPY go-relay-client/ go-relay-client/
RUN go mod download
COPY . .
RUN go get
RUN go build -o ./build/goteth


FROM alpine:latest
# Distrust the emSign ECC Root CA - G3 and AAA Certificate Services roots.
# Alpine ca-certificates 20261020 (NSS 3.131.1) already drops them; the sed
# covers older packages.
RUN apk --no-cache add --upgrade ca-certificates ca-certificates-bundle \
    && sed -i -E 's#^(mozilla/(emSign_ECC_Root_CA_-_G3|Comodo_AAA_Services_root)\.crt)$#!\1#' /etc/ca-certificates.conf \
    && update-ca-certificates
WORKDIR /
COPY --from=builder /app/build/goteth ./
COPY --from=builder /app/pkg/db/migrations ./pkg/db/migrations
ENTRYPOINT ["/goteth"]
