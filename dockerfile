FROM golang:alpine

WORKDIR /app

COPY . /app

RUN go build -o main main.go

EXPOSE 8080

CMD ["go", "run", "main.go"]