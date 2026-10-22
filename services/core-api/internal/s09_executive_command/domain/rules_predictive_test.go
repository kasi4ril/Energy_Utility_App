// Package domain_test verifies the 4 predictive maintenance evaluation models.
package domain_test

import (
	"strings"
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
)

func TestEvaluateHydroWaterAvailability(t *testing.T) {
	t.Run("stable hydro conditions", func(t *testing.T) {
		head, score, level, mwh, summary := domain.EvaluateHydroWaterAvailability(35.0, 30.0, 45.0, 0.90, 14)
		if head <= 45.0 {
			t.Fatalf("expected head to increase or stay steady, got %f", head)
		}
		if score > 0.25 {
			t.Fatalf("expected low score, got %f", score)
		}
		if level != domain.RiskLow {
			t.Fatalf("expected RiskLow, got %s", level)
		}
		if mwh <= 0.0 {
			t.Fatalf("expected positive MWh, got %f", mwh)
		}
		if !strings.Contains(summary, "Stable hydro forecast") {
			t.Fatalf("unexpected summary: %s", summary)
		}
	})

	t.Run("critical low water alert", func(t *testing.T) {
		head, score, level, _, summary := domain.EvaluateHydroWaterAvailability(5.0, 40.0, 15.0, 0.88, 7)
		if level != domain.RiskCritical {
			t.Fatalf("expected RiskCritical, got %s", level)
		}
		if score < 0.75 {
			t.Fatalf("expected critical score, got %f", score)
		}
		if !strings.Contains(summary, "Critical Low Water Alert") {
			t.Fatalf("expected critical summary, got %s", summary)
		}
		_ = head
	})

	t.Run("depletion warning", func(t *testing.T) {
		_, score, level, _, summary := domain.EvaluateHydroWaterAvailability(10.0, 30.0, 21.0, 0.88, 7)
		if level != domain.RiskHigh {
			t.Fatalf("expected RiskHigh, got %s", level)
		}
		if score < 0.50 {
			t.Fatalf("expected high score, got %f", score)
		}
		if !strings.Contains(summary, "Hydro Depletion Warning") {
			t.Fatalf("expected depletion summary, got %s", summary)
		}
	})

	t.Run("moderate drawdown", func(t *testing.T) {
		_, score, level, _, summary := domain.EvaluateHydroWaterAvailability(25.0, 30.0, 35.0, 0.88, 7)
		if level != domain.RiskMedium {
			t.Fatalf("expected RiskMedium, got %s", level)
		}
		if score < 0.25 || score > 0.50 {
			t.Fatalf("expected medium score, got %f", score)
		}
		if !strings.Contains(summary, "Moderate hydro drawdown") {
			t.Fatalf("expected drawdown summary, got %s", summary)
		}
	})
}

func TestEvaluateGeneratorHealth(t *testing.T) {
	t.Run("normal generator parameters", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateGeneratorHealth(2.1, 75.0, 3.2, 2000.0)
		if score != 0.0 {
			t.Fatalf("expected score 0.0, got %f", score)
		}
		if level != domain.RiskLow {
			t.Fatalf("expected RiskLow, got %s", level)
		}
		if autoWO {
			t.Fatal("expected autoWO to be false")
		}
	})

	t.Run("critical vibration and thermal trip", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateGeneratorHealth(12.5, 125.0, 1.0, 9000.0)
		if score < 0.75 {
			t.Fatalf("expected critical score >= 0.75, got %f", score)
		}
		if level != domain.RiskCritical {
			t.Fatalf("expected RiskCritical, got %s", level)
		}
		if !autoWO {
			t.Fatal("expected autoWO to be true for critical alert")
		}
	})

	t.Run("preventive mechanical intervention (high risk)", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateGeneratorHealth(7.5, 106.0, 2.8, 4000.0)
		if score < 0.50 || score >= 0.75 {
			t.Fatalf("expected score in [0.50, 0.75), got %f", score)
		}
		if level != domain.RiskHigh {
			t.Fatalf("expected RiskHigh, got %s", level)
		}
		if !autoWO {
			t.Fatal("expected autoWO to be true for high risk")
		}
	})

	t.Run("medium risk enhanced monitoring", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateGeneratorHealth(4.8, 92.0, 1.6, 1000.0)
		if level != domain.RiskMedium && level != domain.RiskHigh {
			t.Fatalf("expected RiskMedium or RiskHigh, got %s", level)
		}
		_ = autoWO
		_ = score
	})
}

func TestEvaluateTransformerCBM(t *testing.T) {
	t.Run("normal transformer operation", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateTransformerCBM(65.0, 60.0, 0)
		if score != 0.0 {
			t.Fatalf("expected 0 score, got %f", score)
		}
		if level != domain.RiskLow {
			t.Fatalf("expected RiskLow, got %s", level)
		}
		if autoWO {
			t.Fatal("expected autoWO false")
		}
	})

	t.Run("critical transformer overload and top-oil heating", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateTransformerCBM(125.0, 98.0, 5)
		if score < 0.70 {
			t.Fatalf("expected score >= 0.70, got %f", score)
		}
		if level != domain.RiskCritical {
			t.Fatalf("expected RiskCritical, got %s", level)
		}
		if !autoWO {
			t.Fatal("expected autoWO true for critical")
		}
	})

	t.Run("high risk dielectric BDV test recommended", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateTransformerCBM(105.0, 88.0, 1)
		if score < 0.45 {
			t.Fatalf("expected score >= 0.45, got %f", score)
		}
		if level != domain.RiskHigh {
			t.Fatalf("expected RiskHigh, got %s", level)
		}
		if !autoWO {
			t.Fatal("expected autoWO true for high risk")
		}
	})

	t.Run("medium risk loading envelope", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateTransformerCBM(88.0, 70.0, 0)
		if level != domain.RiskMedium && level != domain.RiskLow {
			t.Fatalf("unexpected risk level: %s", level)
		}
		if autoWO {
			t.Fatal("expected autoWO false for medium/low")
		}
		_ = score
	})
}

func TestEvaluateFeederPatrolRisk(t *testing.T) {
	t.Run("normal feeder reliability", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateFeederPatrolRisk(1, 0.3, 30.0, false)
		if score != 0.0 {
			t.Fatalf("expected 0.0 score, got %f", score)
		}
		if level != domain.RiskLow {
			t.Fatalf("expected RiskLow, got %s", level)
		}
		if autoWO {
			t.Fatal("expected autoWO false")
		}
	})

	t.Run("critical feeder bush clearing and insulator wash", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateFeederPatrolRisk(9, 2.5, 120.0, true)
		if score < 0.70 {
			t.Fatalf("expected score >= 0.70, got %f", score)
		}
		if level != domain.RiskCritical {
			t.Fatalf("expected RiskCritical, got %s", level)
		}
		if !autoWO {
			t.Fatal("expected autoWO true")
		}
	})

	t.Run("high risk foot patrol recommended", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateFeederPatrolRisk(5, 1.2, 95.0, false)
		if score < 0.45 {
			t.Fatalf("expected score >= 0.45, got %f", score)
		}
		if level != domain.RiskHigh {
			t.Fatalf("expected RiskHigh, got %s", level)
		}
		if !autoWO {
			t.Fatal("expected autoWO true")
		}
	})

	t.Run("medium risk routine patrol", func(t *testing.T) {
		score, level, _, autoWO := domain.EvaluateFeederPatrolRisk(2, 0.5, 65.0, false)
		if level != domain.RiskMedium && level != domain.RiskLow {
			t.Fatalf("unexpected level: %s", level)
		}
		_ = autoWO
		_ = score
	})
}
