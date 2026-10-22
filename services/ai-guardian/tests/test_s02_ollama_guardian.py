"""Unit and integration test suite for S02 Ollama AI Guardian & Watchdog.

Verifies:
- Maker-Checker fraud risk scoring algorithms & threshold boundaries
- Haversine great-circle GPS proximity & Anti-Ghost detection
- @GridAssistant chat co-pilot store inventory & WorkNotes summarization
- ReportFormatter shift handover markdown generation
- Air-gapped fallback mechanisms when local Ollama daemon is unreachable
"""

import json
import unittest
from unittest.mock import MagicMock, patch

import sys
import os
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from app.core.ollama_client import OllamaClient
from app.agents.fraud_and_gps_watchdog import (
    FraudAndGPSWatchdog,
    calculate_haversine_distance,
)
from app.agents.grid_assistant_chat import GridAssistantChatAgent
from app.agents.report_formatter import ReportFormatter


class TestOllamaClient(unittest.TestCase):
    def test_fallback_heuristic_on_connection_error(self):
        client = OllamaClient(base_url="http://127.0.0.1:9999", model="llama3.2")
        res = client.generate("Please assess risk score for tariff band change")
        self.assertIn("response", res)
        self.assertEqual(res["source"], "heuristic_fallback")
        self.assertIn("Heuristic Audit", res["response"])

    def test_stock_fallback_heuristic(self):
        client = OllamaClient(base_url="http://127.0.0.1:9999", model="llama3.2")
        res = client.generate("Is there stock for meters?")
        self.assertIn("Heuristic Catalog", res["response"])


class TestFraudAndGPSWatchdog(unittest.TestCase):
    def test_tariff_band_downgrade_high_risk(self):
        payload = {"old_band": "A", "new_band": "C"}
        justification = "Customer requested lower minimum hours due to night closure."
        score, reasoning = FraudAndGPSWatchdog.score_maker_checker_risk(
            "TARIFF_BAND_CHANGE", payload, justification
        )
        self.assertGreaterEqual(score, 0.80)
        self.assertIn("High revenue loss alert", reasoning)

    def test_debt_adjustment_boundaries(self):
        # Major debt write-off -> High risk
        score_high, reasoning_high = FraudAndGPSWatchdog.score_maker_checker_risk(
            "DEBT_ADJUSTMENT", {"amount": 75000.0}, "Official waiver approved by management"
        )
        self.assertGreaterEqual(score_high, 0.85)

        # Minor debt write-off -> Low risk
        score_low, _ = FraudAndGPSWatchdog.score_maker_checker_risk(
            "DEBT_ADJUSTMENT", {"amount": 5000.0}, "Rounding correction on ledger reconciliation"
        )
        self.assertLessEqual(score_low, 0.50)

    def test_meter_swap_and_store_stock(self):
        score_meter, _ = FraudAndGPSWatchdog.score_maker_checker_risk(
            "METER_SWAP", {"meter_id": "MTR-01"}, "Burnt terminal block requiring new meter unit"
        )
        self.assertGreaterEqual(score_meter, 0.65)

        score_stock, _ = FraudAndGPSWatchdog.score_maker_checker_risk(
            "STORE_STOCK_EDIT", {"quantity_delta": -20}, "Physical store stock count variance reconciliation"
        )
        self.assertGreaterEqual(score_stock, 0.75)

    def test_short_justification_penalty(self):
        # Very short justification (< 15 chars) incurs penalty
        score_penalized, reasoning = FraudAndGPSWatchdog.score_maker_checker_risk(
            "DEBT_ADJUSTMENT", {"amount": 5000.0}, "Fix balance"
        )
        self.assertIn("Insufficient maker justification", reasoning)

    def test_gps_haversine_distance_and_tolerance(self):
        # Asset located at Jos Substation (approx: 9.9281, 8.8921)
        asset_lat, asset_lon = 9.9281, 8.8921

        # Coordinate 50 meters away
        near_lat, near_lon = 9.9284, 8.8923
        valid, dist, note = FraudAndGPSWatchdog.verify_gps_proximity(
            near_lat, near_lon, asset_lat, asset_lon, tolerance_meters=150.0
        )
        self.assertTrue(valid)
        self.assertLess(dist, 150.0)
        self.assertIn("GPS verified", note)

        # Coordinate 3 kilometers away -> Anti-Ghost alert
        far_lat, far_lon = 9.9500, 8.9200
        valid_far, dist_far, note_far = FraudAndGPSWatchdog.verify_gps_proximity(
            far_lat, far_lon, asset_lat, asset_lon, tolerance_meters=150.0
        )
        self.assertFalse(valid_far)
        self.assertGreater(dist_far, 1000.0)
        self.assertIn("Anti-Ghost Alert", note_far)


class TestGridAssistantChatAgent(unittest.TestCase):
    def setUp(self):
        self.agent = GridAssistantChatAgent()

    def test_store_inventory_lookup(self):
        res = self.agent.process_message(
            "TEAM_LEAD_FEEDER",
            "@GridAssistant do we have any conductors available in store?",
        )
        self.assertEqual(res["intent"], "STORE_STOCK_LOOKUP")
        self.assertIn("Store Inventory Lookup", res["response"])
        self.assertIn("Conductor", res["response"])

    def test_feeder_status_intent(self):
        res = self.agent.process_message(
            "CUSTOMER_CARE_OFFICER",
            "What is the status of the 11kV feeder outage?",
        )
        self.assertEqual(res["intent"], "FEEDER_STATUS")
        self.assertIn("Central Dispatch Hub", res["response"])

    def test_worknotes_summary_intent(self):
        context = {
            "work_notes": [
                {"author_staff_id": "Engr. Dan", "note_text": "Breaker isolated on Feeder 1"},
                {"author_staff_id": "Safety Rep", "note_text": "Earthing verified, permit issued"},
            ]
        }
        res = self.agent.process_message(
            "HEAD_SERVICES",
            "@GridAssistant summarize the latest notes",
            context_data=context,
        )
        self.assertEqual(res["intent"], "WORKNOTES_SUMMARY")
        self.assertIn("Executive Thread Summary", res["response"])


class TestReportFormatter(unittest.TestCase):
    def test_shift_handover_report_formatting(self):
        shift_data = {
            "shift_window": "Morning 08:00 - 16:00",
            "station_name": "Kurra Hydro Master Station",
            "outgoing_operator": "Operator John",
            "incoming_operator": "Operator Sarah",
            "plants_telemetry": [
                {
                    "plant_name": "Plant 1 (Kurra)",
                    "state": "IN_SERVICE",
                    "mw_generated": 14.5,
                    "water_level_m": 842.1,
                    "oil_pressure_bar": 4.2,
                },
                {
                    "plant_name": "Plant 2",
                    "state": "STANDBY",
                    "mw_generated": 0.0,
                    "water_level_m": 841.9,
                    "oil_pressure_bar": 3.8,
                },
            ],
            "active_faults": ["Excitation breaker trip on Unit 3"],
            "safety_advisories": ["Turbine pit clearance permit active until 18:00"],
        }
        md = ReportFormatter.format_shift_handover_report(shift_data)
        self.assertIn("# NESCO OPERATIONAL SHIFT HANDOVER REPORT", md)
        self.assertIn("Plant 1 (Kurra)", md)
        self.assertIn("14.50 MW", md)
        self.assertIn("Excitation breaker trip on Unit 3", md)
        self.assertIn("Turbine pit clearance permit active", md)


if __name__ == "__main__":
    unittest.main()

