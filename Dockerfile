FROM ubuntu:18.04

WORKDIR /app

COPY ./output/goldennum /app

EXPOSE 80

ENTRYPOINT [ "./goldennum" ]