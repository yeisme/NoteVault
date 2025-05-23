FROM docker.io/library/golang:1.24-alpine3.21 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GOTIMEOUT=120s

RUN apk update --no-cache && apk add --no-cache \
    tzdata \
    just \
    curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN just build release

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/notevault /app/notevault
COPY ./etc /app/etc

ENTRYPOINT ["/app/notevault"]

CMD ["server", "-f", "etc/notevaultservice.yaml"]
