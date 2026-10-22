"""Maker-Checker Fraud Risk Scorer, Anti-Ghost & GPS Proximity Watchdog.

Enforces air-gapped security governance for high-risk mutations and verifies
field team proof-of-work coordinates against registered grid assets using Haversine calculation.
"""

from __future__ import annotations

import math
from typing import Any, Dict, Tuple


def calculate_haversine_distance(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
    """Computes great-circle distance between two GPS points on Earth in meters."""
    r_earth = 6371000.0  # Earth's mean radius in meters

    phi1 = math.radians(lat1)
    phi2 = math.radians(lat2)
    delta_phi = math.radians(lat2 - lat1)
    delta_lambda = math.radians(lon2 - lon1)

    a = (
        math.sin(delta_phi / 2.0) ** 2
        + math.cos(phi1) * math.cos(phi2) * (math.sin(delta_lambda / 2.0) ** 2)
    )
    c = 2.0 * math.atan2(math.sqrt(a), math.sqrt(1.0 - a))
    return r_earth * c


class FraudAndGPSWatchdog:
    """Evaluates maker-checker operational risks and audits field GPS proximity."""

    @staticmethod
    def score_maker_checker_risk(
        change_category: str,
        payload: Dict[str, Any],
        justification: str,
    ) -> Tuple[float, str]:
        """Calculates a deterministic AI fraud risk score between 0.00 and 1.00 with audit reasoning."""
        score = 0.20
        reasons = []

        category = change_category.upper()

        if category == "TARIFF_BAND_CHANGE":
            old_band = payload.get("old_band", "")
            new_band = payload.get("new_band", "")
            # Downgrading customer tariff band decreases DisCo revenue (e.g. Band A to C)
            if old_band and new_band and new_band > old_band:
                score = 0.85
                reasons.append(f"High revenue loss alert: tariff band downgraded from {old_band} to {new_band}.")
            else:
                score = 0.45
                reasons.append("Tariff band adjustment requested; requires commercial review.")

        elif category == "DEBT_ADJUSTMENT":
            amount = float(payload.get("amount", 0.0))
            if amount > 50000.0:
                score = 0.90
                reasons.append(f"Substantial financial write-off (NGN {amount:,.2f}) exceeding standard threshold.")
            elif amount > 10000.0:
                score = 0.60
                reasons.append(f"Moderate ledger debt adjustment (NGN {amount:,.2f}).")
            else:
                score = 0.35
                reasons.append(f"Minor debt balance adjustment (NGN {amount:,.2f}).")

        elif category == "METER_SWAP":
            score = 0.70
            reasons.append("Meter replacement presents bypass and vending diversion vulnerability.")

        elif category == "STORE_STOCK_EDIT":
            qty_delta = float(payload.get("quantity_delta", 0.0))
            if qty_delta < 0:
                score = 0.80
                reasons.append("Unscheduled store inventory reduction raises shrinkage risk.")
            else:
                score = 0.40
                reasons.append("Store inventory adjustment upward.")

        elif category == "ALLOWANCE_OVERRIDE":
            score = 0.75
            reasons.append("Staff financial allowance manual override requires 3-way audit reconciliation.")

        else:
            score = 0.50
            reasons.append(f"Standard operational change request under {change_category}.")

        # Justification length and quality penalty
        if len(justification.strip()) < 15:
            score += 0.15
            reasons.append("Insufficient maker justification provided (less than 15 characters).")

        # Clamp score between 0.05 and 0.99
        final_score = max(0.05, min(0.99, round(score, 2)))
        reasoning = " | ".join(reasons)
        return final_score, reasoning

    @staticmethod
    def verify_gps_proximity(
        reported_lat: float,
        reported_lon: float,
        asset_lat: float,
        asset_lon: float,
        tolerance_meters: float = 150.0,
    ) -> Tuple[bool, float, str]:
        """Validates that a field team's reported GPS coordinate is within tolerance of the target asset."""
        dist = calculate_haversine_distance(reported_lat, reported_lon, asset_lat, asset_lon)
        dist_rounded = round(dist, 1)

        if dist <= tolerance_meters:
            return True, dist_rounded, f"GPS verified: Team Lead is {dist_rounded}m from asset (tolerance {tolerance_meters}m)."
        return False, dist_rounded, f"Anti-Ghost Alert: Team Lead is {dist_rounded}m away, exceeding tolerance of {tolerance_meters}m."

