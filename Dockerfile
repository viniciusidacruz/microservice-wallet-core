FROM golang:1.26.2-alpine

WORKDIR /app

RUN apk add --no-cache git bash

CMD ["tail", "-f", "/dev/null"]
