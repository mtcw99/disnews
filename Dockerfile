FROM golang:1.16.3 as builder
WORKDIR /go/src/github.com/mtcw99/disnews
RUN go get github.com/mtcw99/disnews
COPY . ./
RUN GOOS=linux go build -ldflags "-linkmode external -extldflags -static"

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/mtcw99/disnews/disnews .
COPY --from=builder /go/src/github.com/mtcw99/disnews/templates templates
CMD ["./disnews"]

