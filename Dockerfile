FROM golang:1.21

WORKDIR /app

ADD . /app
RUN go mod download

RUN go build -o /renti-backend

EXPOSE 8080

CMD [ "/renti-backend" ]
