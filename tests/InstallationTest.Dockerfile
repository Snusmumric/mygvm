FROM alpine:latest

WORKDIR test 

RUN mkdir gotars
COPY data/go1.16.15.linux-amd64.tar.gz ./gotars/go1.16.15.tar.gz 
COPY data/go1.17.11.linux-amd64.tar.gz ./gotars/go1.17.11.tar.gz

COPY data/mygvm ./mygvm 
