env:
	go get github.com/qiniu/api.v7
	go get github.com/valyala/fasthttp

all:
	go build -o qlist main.go

install: all
	@echo
