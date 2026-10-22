"""HTTP Server Entrypoint for Python Local AI Guardian Service (:8001).

Exposes internal REST endpoints consumed synchronously and asynchronously by Go Core API:
- POST /internal/ai/risk-score (Maker-Checker fraud scoring)
- POST /internal/ai/verify-gps (Haversine field GPS validation)
- POST /internal/ai/chat (@GridAssistant in-app co-pilot)
- POST /internal/ai/format-report (Standardized executive shift reports)
- GET /healthz (Service health check)
"""

from __future__ import annotations

import json
import logging
import os
import sys
from http import HTTPStatus
from http.server import BaseHTTPRequestHandler, HTTPServer
from typing import Any, Dict

from app.agents.fraud_and_gps_watchdog import FraudAndGPSWatchdog
from app.agents.grid_assistant_chat import GridAssistantChatAgent
from app.agents.report_formatter import ReportFormatter

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")


class AIGuardianHandler(BaseHTTPRequestHandler):
    """HTTP request dispatcher for internal AI and Watchdog operations."""

    chat_agent = GridAssistantChatAgent()

    def do_GET(self) -> None:
        """Handles health check inspection."""
        if self.path == "/healthz":
            self._send_json(HTTPStatus.OK, {
                "status": "healthy",
                "service": "ai-guardian",
                "sprint": "s02_ollama_guardian",
            })
        else:
            self._send_json(HTTPStatus.NOT_FOUND, {"error": "not_found", "path": self.path})

    def do_POST(self) -> None:
        """Dispatches internal AI processing requests."""
        content_len = int(self.headers.get("Content-Length", 0))
        raw_body = self.rfile.read(content_len).decode("utf-8") if content_len > 0 else "{}"

        try:
            body = json.loads(raw_body)
        except json.JSONDecodeError:
            self._send_json(HTTPStatus.BAD_REQUEST, {"error": "bad_request", "message": "invalid JSON"})
            return

        if self.path == "/internal/ai/risk-score":
            self._handle_risk_score(body)
        elif self.path == "/internal/ai/verify-gps":
            self._handle_verify_gps(body)
        elif self.path == "/internal/ai/chat":
            self._handle_chat(body)
        elif self.path == "/internal/ai/format-report":
            self._handle_format_report(body)
        else:
            self._send_json(HTTPStatus.NOT_FOUND, {"error": "endpoint_not_found", "path": self.path})

    def _handle_risk_score(self, body: Dict[str, Any]) -> None:
        category = body.get("change_category", "")
        payload = body.get("payload", {})
        justification = body.get("justification", "")

        score, reasoning = FraudAndGPSWatchdog.score_maker_checker_risk(category, payload, justification)
        self._send_json(HTTPStatus.OK, {
            "risk_score": score,
            "reasoning": reasoning,
            "status": "evaluated",
        })

    def _handle_verify_gps(self, body: Dict[str, Any]) -> None:
        try:
            r_lat = float(body["reported_lat"])
            r_lon = float(body["reported_lon"])
            a_lat = float(body["asset_lat"])
            a_lon = float(body["asset_lon"])
            tol = float(body.get("tolerance_meters", 150.0))
        except (KeyError, ValueError) as e:
            self._send_json(HTTPStatus.BAD_REQUEST, {"error": "missing_coordinates", "message": str(e)})
            return

        is_valid, dist, reasoning = FraudAndGPSWatchdog.verify_gps_proximity(r_lat, r_lon, a_lat, a_lon, tol)
        self._send_json(HTTPStatus.OK, {
            "is_valid": is_valid,
            "distance_meters": dist,
            "reasoning": reasoning,
        })

    def _handle_chat(self, body: Dict[str, Any]) -> None:
        role = body.get("sender_role", "STAFF")
        message = body.get("message_text", "")
        context = body.get("context_data")

        res = self.chat_agent.process_message(role, message, context)
        self._send_json(HTTPStatus.OK, res)

    def _handle_format_report(self, body: Dict[str, Any]) -> None:
        report_md = ReportFormatter.format_shift_handover_report(body)
        self._send_json(HTTPStatus.OK, {"report_markdown": report_md})

    def _send_json(self, status: HTTPStatus, data: Dict[str, Any]) -> None:
        resp_bytes = json.dumps(data).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(resp_bytes)))
        self.end_headers()
        self.wfile.write(resp_bytes)

    def log_message(self, format: str, *args: Any) -> None:
        # Suppress noisy standard request log lines during testing
        return


def run_server(port: int = 8001) -> None:
    server_address = ("0.0.0.0", port)
    httpd = HTTPServer(server_address, AIGuardianHandler)
    logging.info(f"NESCO AI Guardian listening on port {port}...")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        logging.info("Shutting down AI Guardian service.")
        httpd.server_close()


if __name__ == "__main__":
    port_arg = int(os.getenv("PORT", "8001"))
    run_server(port_arg)

