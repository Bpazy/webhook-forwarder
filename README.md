# webhook-forwarder
[![Build](https://github.com/Bpazy/webhook-forwarder/workflows/Build/badge.svg)](https://github.com/Bpazy/webhook-forwarder/actions/workflows/build.yml)
[![Test](https://github.com/Bpazy/webhook-forwarder/workflows/Test/badge.svg)](https://github.com/Bpazy/webhook-forwarder/actions/workflows/test.yml)
[![Docker](https://github.com/Bpazy/webhook-forwarder/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/Bpazy/webhook-forwarder/actions/workflows/docker-publish.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Bpazy_webhook-forwarder&metric=alert_status)](https://sonarcloud.io/dashboard?id=Bpazy_webhook-forwarder)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bpazy/webhook-forwarder)](https://goreportcard.com/report/github.com/Bpazy/webhook-forwarder)

Forward the webhook request.

![webhook-forwarder (2)](https://user-images.githubusercontent.com/9838749/205377219-5e0db1d2-6975-43c3-8239-1da1388485cf.png)

## Usage
I suppose your webhook request body looks like this:
```json
{
    "alerts":[
        {
            "status":"resolved",
            "labels":{
                "alertname":"325i alert"
            }
        }
    ]
}
```
And your backend wanted:
```json
{
    "body":"Test Bark Server",
    "title":"bleem"
}
```

Now you can use `webhook-forwarder` to receive and redirect and modify the webhook request like this:
```js
// ~/.config/webhook-forwarder/test.js
function convert(origin) {
    alert = origin.alerts[0];
    return {
        target: ["https://api.day.app/asd/"],
        payload: {
            title: alert.labels.alertname,
            body: "",
        }
    }
};
```

Finally your backend will got correct body.

## Docker Deploy
```yaml
version: '3'
services:
  webhook-forwarder:
    image: ghcr.io/bpazy/webhook-forwarder:master
    environment:
      - PORT=:8080
    restart: always
    ports:
      - "8080:8080"
    volumes:
      - /home/ubuntu/webhook-forwarder/templates:/root/.config/webhook-forwarder/templates
```
