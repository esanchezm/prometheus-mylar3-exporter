FROM golang:1.22-alpine

COPY . /work
WORKDIR /work
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o mylar3-exporter

FROM scratch

COPY --from=0 /work/mylar3-exporter /mylar3-exporter

EXPOSE 4040
USER 65534
ENTRYPOINT ["/mylar3-exporter"]
