
build:
	@echo "Building observer daemon"
	go build -o bin/observer observer.go

install:
	@echo "Installing observer daemon"
	go install observer.go
	@echo "Just installed observer daemon!"
