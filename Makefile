PROXY_SERVER = http://localhost:8080
SET_ENDPOINT = ${PROXY_SERVER}/set
FETCH_ENDPOINT = ${PROXY_SERVER}/fetch
QUERY_ENDPOINT = ${PROXY_SERVER}/query

all:	check

default: check

install: 
	go get github.com/gorilla/mux
	go get gopkg.in/resty.v1

clean:
	rm -rf proxy

build: install clean
	go build proxy.go handler.go

run-server:
	python2.7 server/server.py 8081 &
	python2.7 server/server.py 8082 &
	python2.7 server/server.py 8083 &

run-proxy:
	./proxy localhost:8081 localhost:8082 localhost:8083 &

run-all: run-server run-proxy

kill-server:
	pkill -f server.py

kill-proxy:
	pkill -f proxy

kill-all: kill-server kill-proxy

curl-set:
	curl -H 'Content-Type: application/json' -X PUT -d '[{"key":{"encoding": "asd", "data": "k1"}, "value":{"encoding": "string", "data": "d1"}}, {"key":{"encoding": "string", "data": "k2"}, "value":{"encoding": "string", "data": "d2"}}, {"key":{"encoding": "binary", "data": "k3"}, "value":{"encoding": "string", "data": "d3"}}]' ${SET_ENDPOINT}

curl-fetch-get:
	curl ${FETCH_ENDPOINT}

curl-fetch-post:
	curl -H 'Content-Type: application/json' -X POST -d '[{"encoding": "string", "data": "k2"}, {"encoding": "binary", "data": "k3"}]' ${FETCH_ENDPOINT}

curl-fetch: curl-fetch-get curl-fetch-post

curl-query:
	curl -H 'Content-Type: application/json' -X POST -d '[{"encoding": "string", "data": "k2"}, {"encoding": "binary", "data": "k3"}]' ${QUERY_ENDPOINT}

check: build 
	(make run-all && sleep 2 && make curl-set && make curl-fetch && make curl-query)
	make kill-all

dist:
	dir=`basename $$PWD`; cd ..; tar cvf $$dir.tar ./$$dir; gzip $$dir.tar
	