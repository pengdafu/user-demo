FROM go:1.11.13 AS build

RUN mkdir -p /go/src/user && \
    mkdir -p /go/pkg && \
    mkdir -p /go/bin

ENV GOPATH=/go

RUN cd ${GOPATH} && go get "gopkg.in/go-playground/validator.v10" && \

COPY . ${GOPATH}/src/user/

WORKDIR ${GOPATH}/src/user

RUN govendor sync  && \
    mkdir vendor/github.com/go-playground/validator/v10/ && \
    cp -rf ${GOPATH}/src/gopkg.in/go-playground/validator.v10/* vendor/github.com/go-playground/validator/v10/ && \
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o app cmd/main.gorc

FROM alpine:latest as prod

WORKDIR /go/

COPY --from=builder /build/go/app .

CMD ["./app", "--redisAddr", "redis:6379"]