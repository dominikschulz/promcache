build:
	go build

docker:
	go mod vendor
	docker build -t promcache:latest .

upgrade: clean
	go get -u ./...
	go mod tidy

clean:
	rm -rf promcache vendor

test:
	go test --race -v ./...

.PHONY: build docker upgrade clean test
