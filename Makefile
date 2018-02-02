all:
	export GOPATH=$(pwd)":$GOPATH"
	go get github.com/qiniu/api.v7

install: all
	@echo
