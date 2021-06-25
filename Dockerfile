# compile phase
FROM golang:1.15.12-alpine AS buildStage
ENV GO111MODULE=on
WORKDIR /go/src

ADD go.mod /go/src
ADD go.sum /go/src
RUN go mod download

ADD . /go/src
RUN cd /go/src && go build -o main

# package phase
FROM alpine:latest
WORKDIR /app
COPY --from=buildStage /go/src/main /app/
COPY --from=buildStage /go/src/setting /app/setting

ENTRYPOINT ["/app/main"]

