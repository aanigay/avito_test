FROM golang:1.19-alpine as builder

COPY . /api

WORKDIR /api

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -o api ./cmd/main.go

#**# Это контейнер с ранее собранным приложением**

FROM alpine:latest as server

WORKDIR /root/

COPY --from=builder /api/api .

COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

RUN chown root:root api
RUN chmod +x api

EXPOSE 8081

CMD ["./api"]