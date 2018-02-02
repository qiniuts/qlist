all:
	source env.sh
	go get github.com/qiniu/api.v7

install: all
	@echo
