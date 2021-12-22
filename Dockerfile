FROM golang:1.17-alpine

MAINTAINER Adli I. Ifkar <adly.shadowbane@gmail.com>
ENV LANG=en_US.UTF-8
ENV DEBIAN_FRONTEND noninteractive

USER root

#RUN unset GOPATH
#
#RUN pwd

RUN mkdir -p /opt/application

WORKDIR /opt/application

COPY . .

RUN go build -o /opt/microservice-sekolah cmd/api/main.go

WORKDIR /opt
RUN mkdir -p /opt/log && chmod -R 777 /opt/log
RUN rm -rf /opt/applications

CMD ["/opt/microservice-sekolah"]