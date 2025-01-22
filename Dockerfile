FROM golang:1.23

WORKDIR /app

COPY . .

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]