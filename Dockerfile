FROM golang:latest AS development
RUN apt update
RUN apt install ca-certificates
RUN git clone --progress --verbose --depth=1 https://github.com/Bpazy/webhook-forwarder /webhook-forwarder
WORKDIR /webhook-forwarder
RUN go env && make linux-amd64

FROM ubuntu:latest AS production
ENV PORT ":8080"
COPY --from=development /etc/ssl /etc/ssl
COPY --from=development /webhook-forwarder/bin/webhook-forwarder-linux-amd64 /webhook-forwarder/bin/webhook-forwarder
WORKDIR /webhook-forwarder/bin
ENTRYPOINT ["./webhook-forwarder", "serve", "--port", "echo $PORT", "--verbose"]
