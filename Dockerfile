FROM golang:1.20.4

WORKDIR /root/

COPY . ./

RUN go get
RUN go build

CMD ["./shelly-ht-action-prometheus"]
