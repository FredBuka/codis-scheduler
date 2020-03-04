FROM golang:1.13.4 as builder

ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

WORKDIR /go/src/github.com/oarfah/codis-scheduler
ADD go.mod .
ADD go.sum .
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo .

FROM centos:7

WORKDIR /data/service/scheduler
COPY --from=builder /go/src/github.com/oarfah/codis-scheduler .

CMD ["./codis-scheduler", "--server-replicas=2"]
