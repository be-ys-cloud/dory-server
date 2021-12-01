FROM golang:alpine

RUN mkdir /go/src/dory

WORKDIR /go/src/dory

#RUN go get
ENV GO111MODULE=off
RUN apk add --no-cache git

#RUN go mod init
WORKDIR /go/src/dory/src
ENV GOPATH="/go/src/dory"

RUN apk add build-base

RUN go get -t github.com/go-ldap/ldap
RUN go get -t github.com/sirupsen/logrus
RUN go get -t golang.org/x/text/encoding/unicode
RUN go get -t github.com/thanhpk/randstr
RUN go get -t github.com/gorilla/mux

COPY . /go/src/dory

RUN go build -o /go/bin/dory .

FROM golang:alpine
RUN mkdir /app

COPY --from=0 /go/bin/dory /app/dory
COPY templates/ /app/templates

WORKDIR /app
ENTRYPOINT ["./dory"]

EXPOSE 8000