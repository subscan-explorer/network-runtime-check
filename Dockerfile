FROM golang:1.18 as build

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /app/release

ADD . .
RUN go build -o runtime-check cmd/main.go

FROM alpine as prod

RUN mkdir -p /app/bin/

COPY --from=build /app/release/runtime-check /app/bin/runtime-check

WORKDIR /app/

CMD ["bin/runtime-check"]



