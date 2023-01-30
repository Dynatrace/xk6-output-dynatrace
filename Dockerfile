FROM golang:1.17-alpine as builder
WORKDIR $GOPATH/src/go.k6.io/k6
ADD . .
RUN apk --no-cache add git
RUN CGO_ENABLED=0 go install go.k6.io/xk6/cmd/xk6@latest
RUN CGO_ENABLED=0 xk6 build --output /k6  --with github.com/Dynatrace/xk6-output-dynatrace=.

FROM grafana/k6:latest
COPY --from=builder /k6 /usr/bin/k6
