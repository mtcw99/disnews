db-create:
	docker volume create disnews-db

hr-build:
	docker build -f hotreload.Dockerfile -t mtcw99/disnews:latest .

hr-run:
	docker run -p 8080:8080 -v disnews-db:/go/src/github.com/mtcw99/disnews/db -v $(CURDIR):/go/src/github.com/mtcw99/disnews mtcw99/disnews

