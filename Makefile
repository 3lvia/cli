main_package_path = .
binary_name = 3lv
build_dir = /tmp/3lv/bin
package_dir = /tmp/3lv/dist
go_os = $(shell go env GOOS)
go_arch = $(shell go env GOARCH)

## help: Show this help message.
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## test: Run unit tests.
.PHONY: test
test:
	go test -v ./...

## lint: Run linter (golangci-lint).
.PHONY: lint
lint:
	golangci-lint run ./...

## build: Build the binary (tries to guess the OS and architecture).
.PHONY: build
build:
	GOOS=${go_os} GOARCH=${go_arch} go build -o ${build_dir}/${binary_name} ${main_package_path}

## build-linux-amd64: Build the binary for Linux/amd64.
.PHONY: build-linux-amd64
build-linux-amd64: go_os=linux
build-linux-amd64: go_arch=amd64
build-linux-amd64: build

## build-linux-arm64: Build the binary for Linux/arm64.
.PHONY: build-linux-arm64
build-linux-arm64: go_os=linux
build-linux-arm64: go_arch=arm64
build-linux-arm64: build

## build-macos-amd64: Build the binary for macOS/amd64.
.PHONY: build-macos-amd64
build-macos-amd64: go_os=darwin
build-macos-amd64: go_arch=amd64
build-macos-amd64: build

## build-macos-arm64: Build the binary for macOS/arm64.
.PHONY: build-macos-arm64
build-macos-arm64: go_os=darwin
build-macos-amr64: go_arch=arm64
build-macos-arm64: build

## build-windows-amd64: Build the binary for Windows/amd64.
.PHONY: build-windows-amd64
build-windows-amd64: go_os=windows
build-windows-amd64: go_arch=amd64
build-windows-amd64: build

## run: Build and then run the binary.
.PHONY: run
run: build
	${build_dir}/${binary_name}

## package: Build and then package the binary as a tarball (tries to guess the OS and architecture).
.PHONY: package
package: build
	mkdir -p ${package_dir}
	tar -czf ${package_dir}/3lv-${go_os}-${go_arch}.tar.gz LICENSE README.md -C ${build_dir} 3lv
	cd ${package_dir} && md5sum 3lv-${go_os}-${go_arch}.tar.gz > 3lv-${go_os}-${go_arch}.tar.gz.md5

## package-linux-amd64: Build and then package the binary for Linux/amd64.
.PHONY: package-linux-amd64
package-linux-amd64: go_os=linux
package-linux-amd64: go_arch=amd64
package-linux-amd64: package

## package-linux-arm64: Build and then package the binary for Linux/arm64.
.PHONY: package-linux-arm64
package-linux-arm64: go_os=linux
package-linux-arm64: go_arch=arm64
package-linux-arm64: package

## package-macos-amd64: Build and then package the binary for macOS/amd64.
.PHONY: package-macos-amd64
package-macos-amd64: go_os=darwin
package-macos-amd64: go_arch=amd64
package-macos-amd64: package

## package-macos-arm64: Build and then package the binary for macOS/arm64.
.PHONY: package-macos-arm64
package-macos-arm64: go_os=darwin
package-macos-amr64: go_arch=arm64
package-macos-arm64: package

## package-windows-amd64: Build and then package the binary for Windows/amd64.
.PHONY: package-windows-amd64
package-windows-amd64: go_os=windows
package-windows-amd64: go_arch=amd64
package-windows-amd64: package

## install: Build and then install the binary to /usr/local/bin. Requires root. Only works on Linux and macOS (tries to guess the OS and architecture).
.PHONY: install
install: build
	sudo install -Dm755 -t /usr/local/bin ${build_dir}/${binary_name}

## install-linux-amd64: Build and then install the binary for Linux/amd64 to /usr/local/bin. Requires root.
.PHONY: install-linux-amd64
install-linux-amd64: go_os=linux
install-linux-amd64: go_arch=amd64
install-linux-amd64: install

## install-linux-arm64: Build and then install the binary for Linux/arm64 to /usr/local/bin. Requires root.
.PHONY: install-linux-arm64
install-linux-arm64: go_os=linux
install-linux-arm64: go_arch=arm64
install-linux-arm64: install

## install-macos-amd64: Build and then install the binary for macOS/amd64 to /usr/local/bin. Requires root.
.PHONY: install-macos-amd64
install-macos-amd64: go_os=darwin
install-macos-amd64: go_arch=amd64
install-macos-amd64: install

## install-macos-arm64: Build and then install the binary for macOS/arm64 to /usr/local/bin. Requires root.
.PHONY: install-macos-arm64
install-macos-arm64: go_os=darwin
install-macos-amr64: go_arch=arm64
install-macos-arm64: install
