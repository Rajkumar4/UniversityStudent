FROM Golang:letest
RUN mkdir /build
WORKDIR /build
RUN export GO111MODULE=on
RUN go get github.com/UnivertsityStudent
RUN cd /build && git clone