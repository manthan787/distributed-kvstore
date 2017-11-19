all:	check

default: check

install: 
	go get github.com/gorilla/mux
	go get gopkg.in/resty.v1

build: install
	go build proxy.go handler.go

check: build
	python2.7 server/server.py 8081 &
	python2.7 server/server.py 8082 &
	python2.7 server/server.py 8083 &
	./proxy handler.go localhost:8081 localhost:8082 localhost:8083 &
	curl -H 'Content-Type: application/json' -X PUT -d '[{"key":{"encoding": "asd", "data": "k1"}, "value":{"encoding": "string", "data": "d1"}}, {"key":{"encoding": "string", "data": "k2"}, "value":{"encoding": "string", "data": "d2"}}, {"key":{"encoding": "binary", "data": "k3"}, "value":{"encoding": "string", "data": "d3"}}]' http://localhost:8080/set
	sleep 2
	(ps aux | grep server.py | awk '{print $$2}' | xargs pkill)
	(ps aux | grep proxy | awk '{print $$2}' | xargs pkill)