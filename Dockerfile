FROM ubuntu:18.04

WORKDIR /app

COPY ./output/* /app

EXPOSE 8080

ENTRYPOINT [ "./goldennum" ]