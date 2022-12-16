FROM golang:1.19 as builder
ENV GOPROXY=https://proxy.golang.com.cn,direct
WORKDIR /app
ADD . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o oauth2-server

FROM scratch as final
LABEL maintainer="shengbox@gmail.com"

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /go
COPY --from=builder /app/oauth2-server .
COPY --from=builder /app/web ./web

ENV GIN_MODE=release
ENTRYPOINT ["./oauth2-server"]