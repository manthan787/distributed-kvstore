from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
import urlparse, json

'''
Class for handling requests coming to the server.
'''
class RequestHandler(BaseHTTPRequestHandler):

	def do_GET(self):
		pass

	def do_PUT(self):
		pass
	pass

def start_server(addr, port, handler, verbose=False):
	'''Starts a server at given address and port, with specified `handler`.
	If `verbose` is true, prints extra information on console
	'''
	server = HTTPServer((addr, port), handler)
	if verbose:
		print "Server started at: {}:{}".format(addr, port)
	server.serve_forever()

if __name__ == '__main__':
	import sys
	args = sys.argv
	port = 3000 # Default port for server
	if len(args) == 2:
		port = int(args[1])
	start_server("", port, RequestHandler, verbose=True)