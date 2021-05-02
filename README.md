# disnews
Link aggregation and discussion site written in Go

## License
disnews is released under a the [GNU Affero General Public License v3.0](https://www.gnu.org/licenses/agpl-3.0.html) a free software copyleft license.

## Instructions
You can build and run this program on either bare metal or in a docker container.
The Dockerfile is already setup, so just follow the Docker instructions to get it
up and running.

### Bare-metal
* [Go](https://golang.org/) - Prerequisite (Golang language and tools)
* `go build` - Build the program
* `go run .` - Run the program directly
* `mkdir db` - Make `db` directory for the database file to store in

### Docker
* [Docker](https://www.docker.com/) - Prerequisite (Docker container)
* `docker build -t mtcw99/disnews:latest .` - Build the container
* `docker run -p 8080:8080 mtcw99/disnews` - Run the container

#### Persistent Database
* `docker volume create disnews-db` - Create volume
* `docker run -p 8080:8080 -v disnews-db:/root/db mtcw99/disnews` - Run with the specified volume

#### Hot Reload version
* `docker build -f hotreload.Dockerfile -t mtcw99/disnews:latest .` - Build but with `hotreload.Dockerfile`
* `docker run -p 8080:8080 -v disnews-db:/go/src/github.com/mtcw99/disnews/db -v $(pwd):/go/src/github.com/mtcw99/disnews mtcw99/disnews` - Run it

## Development
* Use `go fmt ./...` before submission 

