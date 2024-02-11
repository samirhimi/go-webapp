FROM golang:1.17-alpine3.14 

WORKDIR /app

COPY main.go go.mod go.sum ./app/

RUN go mod download && \
    go mod verify && \
    go build -o bookstore .

CMD [ "./bookstore" ]



EXPOSE 8080
