FROM golang:1.11-alpine 

MAINTAINER "sam850118sam@gmail.com"

WORKDIR /app

ADD . /app

RUN go get -u github.com/gin-gonic/gin
RUN go get github.com/gomodule/redigo/redis
RUN go get github.com/joho/godotenv

EXPOSE 8080
