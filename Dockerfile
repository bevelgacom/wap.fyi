FROM golang:1.24-alpine as build

RUN apk add --no-cache gcc musl-dev 

ENV CGO_ENABLED=1

COPY ./ /go/src/github.com/bevelgacom/wap.fyi

WORKDIR /go/src/github.com/bevelgacom/wap.fyi

RUN go build -o server ./

FROM alpine:edge

RUN apk add --no-cache gcc musl-dev 

RUN mkdir /opt/wap.fyi
WORKDIR /opt/wap.fyi

COPY --from=build /go/src/github.com/bevelgacom/wap.fyi/stations.csv /opt/wap.fyi/
COPY --from=build /go/src/github.com/bevelgacom/wap.fyi/server /usr/local/bin
COPY --from=build /go/src/github.com/bevelgacom/wap.fyi/templates /opt/wap.fyi/templates

ENTRYPOINT [ "server" ]
