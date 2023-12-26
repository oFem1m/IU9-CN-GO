#!/usr/bin/python3

import cgi
from http.server import BaseHTTPRequestHandler, HTTPServer


class WebServerHandler(BaseHTTPRequestHandler):
        def do_GET(self):

                if self.path.endswith("/abc.html"):
                        self.send_response(200)
                        self.send_header('Content-type', 'text/html; charset=utf-8')
                        self.end_headers()
                        output = """
                        <html>
                        <head>
                        <title>Example</title>
                        <meta charset="utf8">
                        </head>
                        <body>
                        <b>helo DTS21!!!!!!!!</b>
                        </body>
                        </html>
                        """
                        self.wfile.write(output.encode("utf8"))
                        return

port = 8000
server = HTTPServer(('', port), WebServerHandler)
print("Web Server running on port: "+str(port))
server.serve_forever()
