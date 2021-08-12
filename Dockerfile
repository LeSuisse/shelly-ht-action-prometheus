FROM golang:1.16.7-alpine AS build

WORKDIR /go/src/shelly-ht-action-prometheus
COPY . /go/src/shelly-ht-action-prometheus

ENV CGO_ENABLED=0
RUN go build -trimpath

FROM gcr.io/distroless/static-debian10:nonroot

ENV ADDRESS_METRICS 0.0.0.0:17796
ENV ADDRESS_SENSOR 0.0.0.0:17795

COPY --from=build /go/src/shelly-ht-action-prometheus /
CMD ["/shelly-ht-action-prometheus"]
