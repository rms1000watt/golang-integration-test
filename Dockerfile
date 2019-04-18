FROM golang:1.12-alpine AS build_base
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /go/src/github.com/rms1000watt/golang-integration-test
COPY . .
RUN apk add ca-certificates git \
    && go mod download \
    && go build -a -installsuffix cgo -ldflags="-w -s" -o  "/go/bin/person-service"

FROM scratch
COPY --from=build_base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build_base /go/bin/person-service /person-service
ENTRYPOINT [ "/person-service" ]
