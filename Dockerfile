FROM golang:1.21

EXPOSE 8080

WORKDIR /usr/local/app

COPY . .

RUN go build main.go

CMD ["/usr/local/app/main"]
