FROM golang:latest

WORKDIR /forum.github.io

COPY . .

RUN go build -o main .

EXPOSE 8282

CMD ["./main"]
