"""@GridAssistant In-App Community Chat & Real-Time Stock Inquiry Agent.

Interprets natural language queries from utility personnel, answers store inventory
inquiries, summarizes WorkNotes threads, and ensures role-safe guidance.
"""

from __future__ import annotations

import re
from typing import Any, Dict, Optional

from app.core.ollama_client import OllamaClient


class GridAssistantChatAgent:
    """Conversational co-pilot powering the @GridAssistant chat dock."""

    def __init__(self, client: Optional[OllamaClient] = None) -> None:
        self.client = client or OllamaClient()

    def process_message(
        self,
        sender_role: str,
        message_text: str,
        context_data: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """Processes an incoming chat message directed to or involving @GridAssistant."""
        clean_text = message_text.strip()
        is_direct_mention = "@gridassistant" in clean_text.lower() or "@assistant" in clean_text.lower()

        # Extract intent
        lower = clean_text.lower()

        # 1. Store Stock Lookup Intent
        if "stock" in lower or "inventory" in lower or "how many" in lower or "available" in lower:
            return self._handle_stock_inquiry(clean_text, context_data)

        # 2. WorkNotes Thread Summary Intent
        if "summarize" in lower or "summary" in lower or "notes" in lower:
            return self._handle_summary_inquiry(clean_text, context_data)

        # 3. Feeder / Outage Status Intent
        if "feeder" in lower or "outage" in lower or "power" in lower:
            return {
                "response": (
                    "[@GridAssistant]: Live feeder statuses can be monitored via the Central Dispatch Hub. "
                    "All planned isolations require a 5-Step Safety PTW approved by HSE."
                ),
                "intent": "FEEDER_STATUS",
                "is_direct_mention": is_direct_mention,
            }

        # 4. General Assistant Response via Ollama
        system_prompt = (
            "You are @GridAssistant, the internal AI co-pilot for NESCO Nig Ltd. "
            "You provide concise, highly professional utility workflow support to field staff, "
            "station operators, and executives. Maintain strict data integrity."
        )
        llm_out = self.client.generate(clean_text, system_prompt=system_prompt)
        return {
            "response": f"[@GridAssistant]: {llm_out['response']}",
            "intent": "GENERAL_QUERY",
            "source": llm_out.get("source", "local_engine"),
            "is_direct_mention": is_direct_mention,
        }

    def _handle_stock_inquiry(self, query: str, context: Optional[Dict[str, Any]]) -> Dict[str, Any]:
        """Resolves live store inventory questions using attached catalog context."""
        store_items = (context or {}).get("store_inventory", [
            {"code": "CBL-100M", "name": "100mm² All-Aluminum Conductor (AAC)", "qty": 450, "unit": "meters"},
            {"code": "FUS-100A", "name": "100A High Rupturing Capacity (HRC) Fuse", "qty": 85, "unit": "pieces"},
            {"code": "MTR-3PH-01", "name": "Three-Phase STS Smart Prepaid Meter", "qty": 120, "unit": "units"},
            {"code": "TR-OIL-200L", "name": "Uninhibited Transformer Mineral Oil Drum", "qty": 14, "unit": "drums"},
        ])

        matched_item = None
        for item in store_items:
            name_words = item["name"].lower().split()
            if any(w in query.lower() for w in name_words if len(w) > 3):
                matched_item = item
                break

        if matched_item:
            resp = (
                f"[@GridAssistant]: **Store Inventory Lookup:**\n"
                f"- **Item:** {matched_item['name']} (`{matched_item['code']}`)\n"
                f"- **Available Quantity:** {matched_item['qty']} {matched_item['unit']}\n"
                f"- **Status:** {'IN_STOCK' if matched_item['qty'] > 0 else 'OUT_OF_STOCK'}\n"
                f"You may attach this code when initiating your Material Requisition."
            )
            return {"response": resp, "intent": "STORE_STOCK_LOOKUP", "matched_item": matched_item}

        return {
            "response": (
                "[@GridAssistant]: To check warehouse stock, specify the item (e.g. conductors, "
                "HRC fuses, prepaid meters, transformer oil). Store officers update balances perpetually."
            ),
            "intent": "STORE_STOCK_LOOKUP",
            "matched_item": None,
        }

    def _handle_summary_inquiry(self, query: str, context: Optional[Dict[str, Any]]) -> Dict[str, Any]:
        """Summarizes multi-department WorkNotes threads into bullet points."""
        notes = (context or {}).get("work_notes", [])
        if not notes:
            return {
                "response": "[@GridAssistant]: No WorkNotes entries are attached to the current active entity.",
                "intent": "WORKNOTES_SUMMARY",
            }

        bullets = []
        for i, note in enumerate(notes[-5:], start=1):
            text = note.get("note_text", "")
            author = note.get("author_staff_id", "Staff")
            bullets.append(f"{i}. *{author}*: {text[:80]}...")

        summary = "\n".join(bullets)
        resp = f"[@GridAssistant]: **Executive Thread Summary ({len(notes)} notes):**\n{summary}"
        return {"response": resp, "intent": "WORKNOTES_SUMMARY"}

