# Multi stage build
# backend build
FROM golang:1.12-alpine as gobuilder

WORKDIR /app
COPY . ./

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ARG PRODUCT_VERSION="1.0.0"

WORKDIR /app/cmd/imgproc
RUN go build -ldflags "-X 'main.Version=${PRODUCT_VERSION}'" -a -installsuffix cgo -o imgproc .


# result image
FROM alpine:latest

RUN set -eux; \
    apk add --no-cache --virtual ca-certificates && \
    apk add --no-cache  tzdata

COPY --from=gobuilder /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=gobuilder /app/cmd/imgproc/imgproc /imgproc/

COPY init.sh /init.sh
RUN chmod +x /imgproc/imgproc && chmod +x /init.sh

#http port
EXPOSE 80

WORKDIR /mgls

ENTRYPOINT ["/init.sh"]

CMD ["./mgls"]
