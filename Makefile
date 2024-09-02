main_package_path = .
binary_name = 3lv
build_dir = /tmp/3lv/bin
package_dir = /tmp/3lv/dist

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

## build: Build the binary
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ${build_dir}/${binary_name} ${main_package_path}

## run: Build and then run the binary.
.PHONY: run
run: build
	${build_dir}/${binary_name}

## package: Build and then package the binary as a tarball.
.PHONY: build
package: build
	mkdir -p ${package_dir}
	tar -czf ${package_dir}/3lv-linux-amd64.tar.gz LICENSE README.md -C ${build_dir} 3lv
	cd ${package_dir} && md5sum 3lv-linux-amd64.tar.gz > 3lv-linux-amd64.tar.gz.md5

## install: Build and then install the binary to /usr/local/bin. Requires root.
.PHONY: install
install: build
	sudo install -Dm755 -t /usr/local/bin ${build_dir}/${binary_name}
