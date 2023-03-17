export GO111MODULE = on

build: go.sum
	mkdir -p build
	go build -mod=readonly -o build/vccli ./vccli

install: go.sum
	go install -mod=readonly ./vccli

clean:
	rm -rf build
	go clean
