FROM golang:1.14 AS builder
MAINTAINER kel.robert@gmail.com

RUN apt-get update && apt-get install -y \
    fswatch \
    psmisc

WORKDIR /src

COPY . .

RUN make migrate build

FROM ubuntu:bionic AS final

RUN apt-get update && apt-get install --yes --no-install-recommends ca-certificates

COPY --from=builder /src/out/consumer /consumer
COPY ./migrations ./migrations

CMD ["./consumer"]
