env:
	go get github.com/qiniu/api.v7

all:
	go build -o qlist main.go

install: all
	@echo
