#!/usr/bin/env python3
"""Dashboard server with K8s API proxy (avoids CORS issues)."""
import http.server
import urllib.request
import sys
import os

K8S_PROXY = os.environ.get("K8S_PROXY", "http://localhost:8001")

class Handler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path.startswith("/api/"):
            self.proxy_k8s()
        else:
            super().do_GET()

    def proxy_k8s(self):
        url = K8S_PROXY + self.path
        try:
            req = urllib.request.Request(url)
            with urllib.request.urlopen(req, timeout=5) as resp:
                body = resp.read()
                self.send_response(resp.status)
                self.send_header("Content-Type", "application/json")
                self.send_header("Access-Control-Allow-Origin", "*")
                self.end_headers()
                self.wfile.write(body)
        except Exception as e:
            self.send_response(502)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(f'{{"error":"{e}"}}'.encode())

    def log_message(self, format, *args):
        pass  # quiet

port = int(sys.argv[1]) if len(sys.argv) > 1 else 3000
print(f"Dashboard: http://localhost:{port}  (K8s proxy: {K8S_PROXY})")
http.server.HTTPServer(("", port), Handler).serve_forever()
