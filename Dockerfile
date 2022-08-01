FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app
ADD out /app/
ADD static /app/static
ADD docs /app/docs
ADD init.sh /app
ADD config.yaml.example /app/config.yaml.default
