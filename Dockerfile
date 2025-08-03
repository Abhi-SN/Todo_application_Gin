FROM golang:1.23-alpine

WORKDIR /app

COPY go.* ./

RUN go mod download


COPY *.go ./


RUN go build -o /app/todo-app

EXPOSE 8080

CMD [ "/app/todo-app" ]