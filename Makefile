build:
	go build -o project-t cmd/main.go

run: build
	./project-t
