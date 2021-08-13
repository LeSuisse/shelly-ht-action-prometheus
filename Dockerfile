FROM gcr.io/distroless/static-debian10:nonroot

ENV ADDRESS_METRICS 0.0.0.0:17796
ENV ADDRESS_SENSOR 0.0.0.0:17795

COPY shelly-ht-action-prometheus /
CMD ["/shelly-ht-action-prometheus"]
