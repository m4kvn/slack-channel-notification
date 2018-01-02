FROM golang:1.9.2-alpine3.7
RUN apk update && apk upgrade && \
    apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go-wrapper download
RUN go-wrapper install
CMD ["go-wrapper", "run"]