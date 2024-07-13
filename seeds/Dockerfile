# syntax=docker/dockerfile:1

FROM alpine:latest

RUN apk --update add \
    bash \
    postgresql-client  

WORKDIR /app

COPY ./seed.sh /app/seed.sh
COPY ./ /app/seeds


RUN chmod +x /app/seed.sh
CMD ["/bin/bash", "/app/seed.sh"]
