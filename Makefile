FILE=./*.go
ARCH=amd64
OS=Windwos
build:
	@echo "正在编译$(ARCH)架构$(OS)平台的文件"
	GOARCH=$(ARCH) GOOS=$(OS) go build $(FILE)
run: 
	go run ./learn.go