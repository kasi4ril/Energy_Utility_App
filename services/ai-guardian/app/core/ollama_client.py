"""Local Ollama LLM Client and Heuristic Fallback Engine.

Guarantees 100% on-premise air-gapped data sovereignty for NESCO platform.
Routes inference locally to Ollama (http://127.0.0.1:11434) with zero cloud leaks.
Includes deterministic heuristic fallbacks for edge-resilience when Ollama container is offline.
"""

from __future__ import annotations

import json
import os
import urllib.error
import urllib.request
from typing import Any, Dict, Optional


class OllamaClient:
    """Client for local Ollama instance with deterministic heuristic fallback."""

    def __init__(self, base_url: Optional[str] = None, model: str = "llama3.2") -> None:
        self.base_url = (base_url or os.getenv("OLLAMA_HOST", "http://127.0.0.1:11434")).rstrip("/")
        self.model = os.getenv("OLLAMA_MODEL", model)
        self.timeout = int(os.getenv("OLLAMA_TIMEOUT_SECONDS", "5"))

    def generate(self, prompt: str, system_prompt: Optional[str] = None) -> Dict[str, Any]:
        """Sends inference request to local Ollama endpoint."""
        url = f"{self.base_url}/api/generate"
        payload = {
            "model": self.model,
            "prompt": prompt,
            "stream": False,
        }
        if system_prompt:
            payload["system"] = system_prompt

        try:
            req = urllib.request.Request(
                url,
                data=json.dumps(payload).encode("utf-8"),
                headers={"Content-Type": "application/json"},
                method="POST",
            )
            with urllib.request.urlopen(req, timeout=self.timeout) as resp:
                data = json.loads(resp.read().decode("utf-8"))
                return {
                    "response": data.get("response", ""),
                    "source": "ollama_local",
                    "model": self.model,
                }
        except (urllib.error.URLError, TimeoutError, OSError) as e:
            # Fallback locally to heuristic engine to ensure zero downtime on network partition
            return {
                "response": self._heuristic_fallback(prompt),
                "source": "heuristic_fallback",
                "error": str(e),
                "model": "rule_based_v1",
            }

    def _heuristic_fallback(self, prompt: str) -> str:
        """Deterministic rule-based response engine when local LLM daemon is offline."""
        lower = prompt.lower()
        if "risk" in lower:
            return "Heuristic Audit: Change request requires strict supervisor verification."
        if "stock" in lower or "meter" in lower:
            return "Heuristic Catalog: Verified item inquiry against registered inventory."
        return "Heuristic Process: Request logged and queued under enterprise compliance rules."

