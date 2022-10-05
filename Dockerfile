FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app
ADD out /app/
