FROM golang

WORKDIR mygvm_tests

RUN make build
RUN ./bin/mygvm -src tar install go1.17.11 tests/data/go1.17.11.linux-amd64.tar.gz
