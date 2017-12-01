FROM golang:1.9-alpine
MAINTAINER mateus.bona@gmail.com

ENV SENDGRID_API_KEY YOUR-SENDGRID-API-KEY

WORKDIR /
RUN mkdir app

WORKDIR /app
COPY monitoramento.go .

RUN apk update && apk add git tzdata && cp /usr/share/zoneinfo/Brazil/East /etc/localtime && go get -v github.com/sendgrid/sendgrid-go && go get -v github.com/gorilla/mux && go build monitoramento.go

EXPOSE 8000

ENTRYPOINT /app/monitoramento
