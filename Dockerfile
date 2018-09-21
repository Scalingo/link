FROM golang:1.11.0

RUN go get github.com/cespare/reflex
ADD . /go/src/github.com/Scalingo/link
WORKDIR /go/src/github.com/Scalingo/link
EXPOSE 1313
RUN go install
CMD /go/bin/link

