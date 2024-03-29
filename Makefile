clean-data:
	rm -fr data/*

slave: clean-data
	go run main.go slave

build-docker:
	GOOS="linux" GOARCH="amd64" go build
	docker build -t ds .

.PHONY: slave clean-data build-docker
