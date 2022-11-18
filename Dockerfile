FROM golang:lastest

RUN mkdir/build
WORKDIR /build

RUN export GO111MODULE=on
RUN go get github.com/l3ankZKuma/os_project
RUN cd /build && git clone https://github.com/l3ankZKuma/os_project.git

RUN cd/build/os_project && go build
EXPOSE 8080

ENTRYPOINT ["/build/os_project/main"]