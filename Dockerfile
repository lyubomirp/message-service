FROM golang:1.20.5

WORKDIR /app

COPY . .
RUN go get  -t -v ./...
RUN go mod tidy

COPY *.go ./

RUN go build -o /message-service

EXPOSE 8080

CMD [ "/message-service" ]