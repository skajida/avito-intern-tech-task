FROM golang:alpine3.18 as builder

ENV PACKAGE service-segs
ENV DIRPATH /src/skajida/$PACKAGE
WORKDIR $DIRPATH

# dependencies
COPY go.mod go.sum .
RUN go mod download

# build
COPY . .
ENV CGO_ENABLED 0
ENV GOOS linux
RUN go build -a -o bin/$PACKAGE $DIRPATH/cmd/$PACKAGE

# image
FROM alpine:3.18.3
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /src/skajida/service-segs/bin/service-segs /usr/bin/service-segs
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/service-segs"]
