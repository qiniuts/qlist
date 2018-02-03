all:
	export GOPATH=$(pwd)":$GOPATH"
	go get github.com/qiniu/api.v7
	go build -o qlist main.go

install: all
	@echo
