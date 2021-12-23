FROM golang:1.17 as builder

WORKDIR /go/src/github.com/ez-deploy/cluster
COPY . .

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.io && \
    CGO_ENABLED=0 go build -tags netgo -o cluster ./main.go

FROM busybox

WORKDIR /

COPY --from=builder /go/src/github.com/ez-deploy/cluster/cluster /cluster

EXPOSE 80
ENTRYPOINT [ "/cluster" ]