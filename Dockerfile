# Choose whatever you want, version >= 1.16
FROM golang:1.21-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.40.4

COPY . .
RUN go mod download
EXPOSE 8080

CMD ["air", "-c", ".air.toml"]