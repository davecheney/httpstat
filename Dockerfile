FROM golang:1.9 as build

RUN go get github.com/davecheney/httpstat && \
    CGO_ENABLED=0 GOOS=linux go build github.com/davecheney/httpstat && \
    go test github.com/davecheney/httpstat

FROM scratch
COPY --from=build /go/httpstat /
CMD ["/httpstat"]
