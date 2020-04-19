FROM golang:1.14-alpine AS builder
MAINTAINER kel.robert@gmail.com

WORKDIR /src

COPY . .

RUN apk add --update git make
RUN make build

FROM ubuntu:bionic AS final-image

RUN apt-get update && apt-get install --yes --no-install-recommends ca-certificates

COPY --from=builder /src/out/consumer /consumer

CMD ["./consumer"]
