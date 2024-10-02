FROM golang:1.22

RUN go install github.com/cespare/reflex@latest
ADD . /go/src/github.com/Scalingo/link
WORKDIR /go/src/github.com/Scalingo/link
EXPOSE 1313
RUN go install -buildvcs=false
CMD /go/bin/link

