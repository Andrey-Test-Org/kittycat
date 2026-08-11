"""Tiny Flask-style HTTP server stub for the kitty dashboard."""
from __future__ import annotations
import json
from http.server import BaseHTTPRequestHandler, HTTPServer
from typing import Dict, List

from .aggregator import aggregate
from .models import KittyEvent

EVENTS: List[KittyEvent] = []


class DashboardHandler(BaseHTTPRequestHandler):
    def do_GET(self) -> None:
        if self.path == "/healthz":
            self._json(200, {"ok": True})
            return
        if self.path == "/stats":
            stats = aggregate(EVENTS)
            payload = {cid: s.__dict__ for cid, s in stats.items()}
            self._json(200, payload)
            return
        self._json(404, {"error": "not found"})

    def _json(self, code: int, body: Dict) -> None:
        data = json.dumps(body, default=str).encode()
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)


def serve(host: str = "127.0.0.1", port: int = 8765) -> None:
    HTTPServer((host, port), DashboardHandler).serve_forever()


if __name__ == "__main__":
    serve()
