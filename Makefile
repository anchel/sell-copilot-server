
NAME=sellcopilot
BUILD_DIR=build
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s -buildid='

normal: clean all

clean:
	rm -rf $(BUILD_DIR)/*

all: darwin-arm64 linux-amd64

darwin-arm64:
	mkdir -p $(BUILD_DIR)/$@ 
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BUILD_DIR)/$@/$(NAME)

linux-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BUILD_DIR)/$@/$(NAME)

test:
	go test -v ./...

run:
	mkdir -p dist
	go build -o dist/sellcopilot main.go
	cp .env dist
	cd dist && ./sellcopilot config.toml
