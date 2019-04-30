FROM golang:latest

MAINTAINER "sam850118sam@gmail.com"

WORKDIR /app

ADD . /app

RUN go get -u github.com/gin-gonic/gin
RUN go get github.com/gomodule/redigo/redis

EXPOSE 3001
