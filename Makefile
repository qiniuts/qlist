env:
	go get github.com/qiniu/api.v7

linux:
	GOOS=linux GOARCH=amd64 go build -o qlist_linux main.go

all:
	GOOS=linux GOARCH=amd64 go build -o qlist main.go

install: all
	@echo
