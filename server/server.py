from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
from store import InMemoryStore
import urlparse, json
from response_codes import *

store = InMemoryStore(set(["string", "binary"]))

class RequestHandler(BaseHTTPRequestHandler):
    """
    Class for handling requests coming to the server.

    Server tests:

    GET /fetch:
    curl http://localhost:3000/fetch

    POST /fetch:
    curl -X POST -H "Content-Type: application/json" -d '[{"encoding": "string", "data": "key1"}, {"encoding": "string", "data": "key176"}]' http://localhost:3000/fetch

    GET /query:
    curl http://localhost:3000/query

    PUT /set:
    curl -X PUT -H "Content-Type: application/json" -d '[{"key": {"encoding" : "string", "data": "key1"}, "value": {"encoding" : "string", "data" : "value1"}}, {"key": {"encoding" : "string", "data": "key2"}, "value": {"encoding" : "string", "data" : "value2"}}]' http://localhost:3000/set
    """
    def _set_response_headers(self, response_code):
        #self.send_header("Content-Type", "application/json")
        self.send_response(response_code)
        self.end_headers()

    def do_GET(self):
        """ Handles GET requests """
        print self.headers
        if not self._path_equals("/fetch"): return self.send_error(NOT_FOUND)
        all_kvs = list(store.get_all())
        self._set_response_headers(SUCCESS_CODE)
        self.wfile.write(json.dumps(all_kvs))

    def do_PUT(self):
        """ Handles PUT requests """
        if not self._path_equals("/set"): return self.send_error(NOT_FOUND)
        if not self._is_json_request(): return self.send_error(BAD_REQUEST)
        request_data = self._read_request_body()
        if not request_data: return self.send_error(FORBIDDEN)
        stats = store.batch_put(request_data)
        self._set_response_headers(SUCCESS_CODE)
        self.wfile.write(json.dumps(stats))

    def do_POST(self):
        """ Handles POST requests """
        if not self._is_json_request(): return self.send_error(BAD_REQUEST)
        if self._path_equals("/fetch"):
            keys = self._read_request_body()
            if not keys: return self.send_error(FORBIDDEN)
            self._handle_post_paths(keys, list(store.batch_get(keys)))
        if self._path_equals("/query"):
            keys = self._read_request_body()
            if not keys: return self.send_error(FORBIDDEN)
            self._handle_post_paths(keys, list(store.batch_lookup(keys)))

    def _handle_post_paths(self, keys, result):
        if len(keys) == len(result):
            self._set_response_headers(SUCCESS_CODE)
        else:
            self._set_response_headers(FORBIDDEN)
        return self.wfile.write(json.dumps(result))

    def _is_json_request(self):
        """ Returns true if the request content-type is json, false otherwise """
        return self.headers.get('Content-Type', "") == 'application/json'

    def _path_equals(self, path):
        """ Returns True if request path equals `path`, False otherwise """
        return self.path.lower() == path

    def _read_request_body(self):
        """ Returns the request body as json """
        try: return json.loads(self.rfile.read(int(self.headers['Content-Length'])))
        except: return ""

def start_server(addr, port, handler):
    """ Starts a server at given address and port, with specified `handler`.
    If `verbose` is true, prints extra information on console
    """
    server = HTTPServer((addr, port), handler)
    print "Server started at: {}:{}".format(addr, port)
    server.serve_forever()

if __name__ == '__main__':
    import sys
    args = sys.argv
    port = 3000 # Default port for server
    if len(args) == 2: port = int(args[1])
    start_server("", port, RequestHandler)