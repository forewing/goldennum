FROM ubuntu:18.04

RUN apt-get update -yq && \
    apt-get install dumb-init

WORKDIR /app

COPY output /app/

EXPOSE 8080

ENTRYPOINT [ "dumb-init", "./goldennum" ]