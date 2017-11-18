from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
from store import InMemoryStore
import urlparse, json

store = InMemoryStore()
store.put({"encoding" : "string", "data": "hello"}, {"encoding" : "string", "data" : "world"})
print store.encodings
print store.data
print list(store.get_all())
'''
Class for handling requests coming to the server.

Server tests:

GET /fetch:
curl http://localhost:3000/fetch

GET /query:
curl http://localhost:3000/query

PUT /set:
curl -X PUT -H "Content-Type: application/json" -d '[{"key1":"value1", "key2": "value2"}]' http://localhost:3000/set
'''
class RequestHandler(BaseHTTPRequestHandler):

    def _set_response_headers(self, response_code):
        self.send_header("Content-Type", "application/json")
        self.send_response(response_code)
        self.end_headers()

    def do_GET(self):
        if self._path_equals("/fetch"):
            print "A fetch request has been made"
            all_kvs = list(store.get_all())
            self._set_response_headers(200)
            print all_kvs
            return self.wfile.write(json.dumps(all_kvs))
        if self._path_equals("/query"):
            print "A query request has been made"
            self.wfile.write(self)
        self.send_error(404)

    def do_PUT(self):
        if not self._is_request_valid(): return self.send_error(400)
        print self.rfile.read(int(self.headers['Content-Length']))
        if self._path_equals("/set"): pass
        self.wfile.write(self)

    def do_POST(self):
        if not self._is_request_valid(): return self.send_error(400)
        if self._path_equals("/fetch"):
            print "A fetch request has been made"

        if self._path_equals("/query"):
            print "A query request has been made"

    def _is_request_valid(self):
        return self.headers.get('Content-Type', "") == 'application/json'

    def _path_equals(self, path):
        return self.path.lower() == path

def start_server(addr, port, handler, verbose=False):
    '''Starts a server at given address and port, with specified `handler`.
    If `verbose` is true, prints extra information on console
    '''
    server = HTTPServer((addr, port), handler)
    if verbose: print "Server started at: {}:{}".format(addr, port)
    server.serve_forever()

if __name__ == '__main__':
    import sys
    args = sys.argv
    port = 3000 # Default port for server
    if len(args) == 2: port = int(args[1])
    start_server("", port, RequestHandler, verbose=True)