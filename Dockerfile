FROM golang:1.14-alpine AS builder
RUN apk add --no-cache \
    build-base

WORKDIR /build
COPY . /build/
RUN ./build.sh


FROM alpine:3
RUN apk add --no-cache \
    bash \
    dumb-init

WORKDIR /app
COPY --from=builder /build/output /app/

EXPOSE 8080
ENTRYPOINT [ "dumb-init", "--" ]
CMD ["/app/goldennum"]