FROM alpine:latest

WORKDIR /app
COPY out /app
ADD static /app/static
ADD docs /app/docs
ADD init.sh /app
ADD config.yaml.example /app/config.yaml.default

ENTRYPOINT [ "/app/api", "server" ]