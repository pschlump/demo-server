
all: go-server

setup:
	mkdir -p ./log

go-server: 
	mkdir -p ./log
	go build 

run0:
	rm -f ./demo-server 
	mkdir -p ./log
	go build
	./demo-server -D db-func-call,IsLoggedIn,CheckAuthTokenInDb -l ./log/log.log


run1:
	rm -f ./demo-server 
	mkdir -p ./log
	go build
	echo ./demo-server -D db-func-call,IsLoggedIn,CheckAuthTokenInDb -l ./log/log.log -E
	./demo-server -l ./log/log.log -E


run2:
	rm -f ./demo-server 
	mkdir -p ./log
	go build
	./demo-server -l ./log/log.log
