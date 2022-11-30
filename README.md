# webhook-forwarder
[![Build](https://github.com/Bpazy/webhook-forwarder/workflows/Build/badge.svg)](https://github.com/Bpazy/webhook-forwarder/actions/workflows/build.yml)
[![Test](https://github.com/Bpazy/webhook-forwarder/workflows/Test/badge.svg)](https://github.com/Bpazy/webhook-forwarder/actions/workflows/test.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Bpazy_webhook-forwarder&metric=alert_status)](https://sonarcloud.io/dashboard?id=Bpazy_webhook-forwarder)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bpazy/webhook-forwarder)](https://goreportcard.com/report/github.com/Bpazy/webhook-forwarder)

Forward the webhook request.

Templates:
```js
function convert(origin) {
    return {
        target: "https://www.baidu.com",
        payload: {
            "hello": origin.name,
        }
    }
}
```