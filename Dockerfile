FROM golang:latest AS development
RUN git clone --progress --verbose --depth=1 https://github.com/Bpazy/webhook-forwarder /webhook-forwarder
WORKDIR /webhook-forwarder
RUN go env && make linux-amd64

FROM ubuntu:latest AS production
RUN update-ca-certificates
ENV PORT 8080
COPY --from=development /webhook-forwarder/bin/webhook-forwarder-linux-amd64 /webhook-forwarder/bin/webhook-forwarder
WORKDIR /webhook-forwarder/bin
ENTRYPOINT ./webhook-forwarder serve --port $PORT
