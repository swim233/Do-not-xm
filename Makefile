FILE=./*.go
ARCH=amd64
OS=windows
FILENAME=release
build:
	@echo "正在编译$(ARCH)架构$(OS)平台的文件"
	GOARCH=$(ARCH) GOOS=$(OS) go build $(FILE)
run: 
	go run ./learn.go
build_all:
	@echo "正在编译全平台的文件"
	GOARCH=amd64 GOOS=linux go build -o $(FILENAME)-amd64-linux $(FILE)
	GOARCH=arm64 GOOS=linux go build -o $(FILENAME)-arm64-linux $(FILE)
	GOARCH=amd64 GOOS=windows go build -o $(FILENAME)-amd64-windows $(FILE)
	GOARCH=arm64 GOOS=windows go build -o $(FILENAME)-arm64-windows $(FILE)