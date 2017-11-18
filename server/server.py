from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
from store import InMemoryStore
import urlparse, json
from response_codes import *

store = InMemoryStore()
store.put({"encoding" : "string", "data": "hello"}, {"encoding" : "string", "data" : "world"})
print store.encodings
print store.data
print list(store.get_all())
print store.batch_put([{"key": {"encoding" : "string", "data": "hell"}, "value": {"encoding" : "string", "data" : "world"}}])
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
        ''' Handles GET requests '''
        if self._path_equals("/fetch"):
            print "A fetch request has been made"
            all_kvs = list(store.get_all())
            if all_kvs:
                self._set_response_headers(SUCCESS_CODE)
                return self.wfile.write(json.dumps(all_kvs))
            return self.send_error(BAD_REQUEST)
        if self._path_equals("/query"):
            print "A query request has been made"
            self.wfile.write(self)
        self.send_error(NOT_FOUND)

    def do_PUT(self):
        ''' Handles PUT requests '''
        if not self._is_json_request(): return self.send_error(BAD_REQUEST)
        request = self._read_request_body()
        if self._path_equals("/set"):
            pass
        self.wfile.write(self)

    def do_POST(self):
        ''' Handles POST requests '''
        if not self._is_json_request(): return self.send_error(BAD_REQUEST)
        if self._path_equals("/fetch"):
            print "A fetch request has been made"
        if self._path_equals("/query"):
            print "A query request has been made"

    def _is_json_request(self):
        ''' Returns true if the requrst content type is json, false otherwise'''
        return self.headers.get('Content-Type', "") == 'application/json'

    def _path_equals(self, path):
        ''' Returns True if request path equals `path`, False otherwise '''
        return self.path.lower() == path

    def _read_request_body(self):
        ''' Returns the request body as json'''
        try: return json.loads(self.rfile.read(int(self.headers['Content-Length'])))
        except: return ""

def start_server(addr, port, handler):
    '''Starts a server at given address and port, with specified `handler`.
    If `verbose` is true, prints extra information on console
    '''
    server = HTTPServer((addr, port), handler)
    print "Server started at: {}:{}".format(addr, port)
    server.serve_forever()

if __name__ == '__main__':
    import sys
    args = sys.argv
    port = 3000 # Default port for server
    if len(args) == 2: port = int(args[1])
    start_server("", port, RequestHandler)