FROM golang:alpine

RUN mkdir /go/src/dory

WORKDIR /go/src/dory

#RUN go get
RUN apk add --no-cache git

#RUN go mod init
WORKDIR /go/src/dory

RUN apk add build-base

COPY . /go/src/dory

RUN go mod download
RUN go build -o /go/bin/dory ./cmd

FROM golang:alpine
RUN mkdir /app

COPY --from=0 /go/bin/dory /app/dory
COPY templates/ /app/templates

WORKDIR /app
ENTRYPOINT ["./dory"]

EXPOSE 8000