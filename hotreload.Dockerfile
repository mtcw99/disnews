FROM golang:1.16.3
WORKDIR /go/src/github.com/mtcw99/disnews
COPY . .
RUN go get github.com/githubnemo/CompileDaemon
RUN mkdir db
ENTRYPOINT CompileDaemon -build="go build ." -command="./disnews"

