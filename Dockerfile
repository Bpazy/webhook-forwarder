FROM golang:latest AS development
RUN git clone --progress --verbose --depth=1 https://github.com/Bpazy/webhook-forwarder /webhook-forwarder
WORKDIR /webhook-forwarder
RUN go env && make linux-amd64

FROM ubuntu:latest AS production
COPY entrypoint.sh /usr/local/bin
COPY --from=development /webhook-forwarder/bin/webhook-forwarder-linux-amd64 /usr/local/bin/webhook-forwarder
ENTRYPOINT entrypoint.sh
