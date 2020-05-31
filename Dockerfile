FROM golang:alpine as builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG store=file
RUN go build -tags jsoniter,${store} -i -o build/app cmd/webhookdelia/main.go

FROM alpine
RUN apk --no-cache --update add ca-certificates && update-ca-certificates
WORKDIR /bin
COPY --from=builder /src/build/app /bin/app
ENTRYPOINT app
