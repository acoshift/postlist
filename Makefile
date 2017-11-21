dev:
	go run cmd/postlist/main.go

clean:
	rm -f postlist

build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o postlist -a -ldflags '-w -s' cmd/postlist/main.go

docker: clean build
	docker build -t postlist .

