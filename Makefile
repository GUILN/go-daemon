all: build build_installer copy_config
	@echo "all good to proceed!"
	@echo "in order to continue the installation of the daemon access bin folder and run the install_observer program"
	exit 0
build:
	@echo "Building observer daemon"
	go build -o bin/observer observer.go
	@echo "Just built observer daemon!"

build_installer:
	@echo "Building observer's installer"
	go build -o bin/install_observer installer/installer.go
	@echo "Just built observer's installer"

copy_config:
	@echo "Copying config file"
	cp config.conf bin/observer.conf
	@echo "Just copied config file"

install:
	@echo "Installing observer daemon"
	go install observer.go
	@echo "Just installed observer daemon!"
