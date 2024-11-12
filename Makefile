FILE=./*.go
ARCH=amd64
OS=linux
FILENAME=release
VERSION=1.0
build:
	@echo "正在编译$(ARCH)架构$(OS)平台的文件"
	GOARCH=$(ARCH) GOOS=$(OS) go build -o $(FILENAME)-$(ARCH)-linux-$(VERSION) $(FILE)
	@echo "文件$(FILENAME)-$(ARCH)-linux-$(VERSION) 编译完成"
run: 
	go run ./*.go
build_all:
	@echo "正在编译全平台的文件"
	GOARCH=amd64 GOOS=linux go build -o $(FILENAME)-amd64-linux-$(VERSION) $(FILE)
	GOARCH=arm64 GOOS=linux go build -o $(FILENAME)-arm64-linux-$(VERSION) $(FILE)
	GOARCH=amd64 GOOS=windows go build -o $(FILENAME)-amd64-windows-$(VERSION) $(FILE)
	GOARCH=arm64 GOOS=windows go build -o $(FILENAME)-arm64-windows-$(VERSION) $(FILE)
test:
	@echo "正在编译全平台的文件"
	GOARCH=amd64 GOOS=linux go build -o $(FILENAME)-amd64-linux-$(VERSION)-test $(FILE)
	GOARCH=arm64 GOOS=linux go build -o $(FILENAME)-arm64-linux-$(VERSION)-test $(FILE)
	GOARCH=amd64 GOOS=windows go build -o $(FILENAME)-amd64-windows-$(VERSION)-test $(FILE)
	GOARCH=arm64 GOOS=windows go build -o $(FILENAME)-arm64-windows-$(VERSION)-test $(FILE)
	@sleep 3s
	@rm ./*test
	@echo "构筑文件清理完成"