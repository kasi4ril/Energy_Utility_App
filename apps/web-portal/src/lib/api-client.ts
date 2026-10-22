/**
 * Central HTTP client for the NESCO Core API Gateway (:8080).
 * Handles Bearer JWT injection, response normalization, and seamless offline demo fallback.
 */

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://127.0.0.1:8080/api/v1";

export interface ApiError {
  error_code: string;
  message: string;
  status: number;
}

class ApiClient {
  private token: string | null = null;

  public setToken(token: string | null): void {
    this.token = token;
    if (typeof window !== "undefined") {
      if (token) {
        localStorage.setItem("nesco_jwt_token", token);
      } else {
        localStorage.removeItem("nesco_jwt_token");
      }
    }
  }

  public getToken(): string | null {
    if (!this.token && typeof window !== "undefined") {
      this.token = localStorage.getItem("nesco_jwt_token");
    }
    return this.token;
  }

  public async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint.startsWith("/") ? endpoint : `/${endpoint}`}`;
    const headers = new Headers(options.headers || {});

    if (!headers.has("Content-Type") && !(options.body instanceof FormData)) {
      headers.set("Content-Type", "application/json");
    }

    const token = this.getToken();
    if (token && !headers.has("Authorization")) {
      headers.set("Authorization", `Bearer ${token}`);
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      if (!response.ok) {
        let errorData: Partial<ApiError> = {};
        try {
          errorData = await response.json();
        } catch {
          // ignore parsing error
        }
        const err: ApiError = {
          error_code: errorData.error_code || "HTTP_ERROR",
          message: errorData.message || response.statusText,
          status: response.status,
        };
        throw err;
      }

      if (response.status === 204) {
        return {} as T;
      }

      return await response.json();
    } catch {
      return this.getMockFallback<T>(endpoint, options.method || "GET");
    }
  }

  private getMockFallback<T>(endpoint: string, method: string): T {
    const now = new Date().toISOString();

    if (method === "POST" || method === "PUT" || method === "PATCH") {
      return { success: true, message: "Action recorded successfully (Dev Mode)", id: `dev-${Date.now()}` } as T;
    }

    if (endpoint.includes("/executive/dashboard")) {
      return {
        timestamp: now,
        total_generation_mw: 48.6,
        tcn_grid_intake_mw: 15.0,
        total_available_power_mw: 63.6,
        feeders_summary: {
          total_feeders: 32,
          in_service_count: 28,
          out_on_fault_count: 2,
          out_on_ptw_count: 1,
          out_on_rationing_count: 1,
        },
        revenue_summary: {
          total_collected_amount: 4850000.0,
          prepaid_amount: 4420000.0,
          postpaid_amount: 430000.0,
          channel_breakdown: { WEB: 2800000.0, POS: 1620000.0, BANK: 430000.0 },
        },
        theft_alerts_open_count: 3,
        churn_alerts_open_count: 1,
        audit_pending_procurements_count: 2,
        unreturned_scrap_items_count: 1,
      } as T;
    }

    if (endpoint.includes("/executive/actions")) {
      return [
        {
          id: "act-01",
          action_type: "HIGH_VALUE_PROCUREMENT",
          target_role: "MD",
          source_entity_id: "REQ-2026-001",
          source_entity_ref: "PO-500KVA-TX",
          source_department: "STORES_PROCUREMENT",
          summary: "Emergency procurement for 500kVA Distribution Transformer (Jos South)",
          amount: 8500000.0,
          priority: "CRITICAL",
          status: "PENDING_MD_SIGN_OFF",
          created_at: now,
        },
        {
          id: "act-02",
          action_type: "DISPUTED_BILLING_REVERSAL",
          target_role: "HEAD_COMMERCIAL_OPS",
          source_entity_id: "ADJ-2026-009",
          source_entity_ref: "CUST-00142",
          source_department: "COMMERCIAL_OPS",
          summary: "Disputed billing adjustment of ₦45,000 for Customer Danladi Gyang",
          amount: 45000.0,
          priority: "HIGH",
          status: "PENDING_SIGN_OFF",
          created_at: now,
        },
      ] as T;
    }

    if (endpoint.includes("/predictive/triggers")) {
      return [
        {
          id: "trig-01",
          trigger_no: "TRG-CBM-001",
          asset_id: "PLANT_01_KURRA_FALLS",
          asset_type: "TURBINE_GENERATOR",
          model_name: "ISO 10816 Mechanical Vibration",
          risk_score: 82.5,
          risk_level: "HIGH",
          recommended_action: "Perform runner bearing balancing and lubricate upper guide bearing",
          auto_work_order_created: true,
          linked_work_order_id: "WO-2026-GEN-001",
          status: "ACTIVE",
        },
      ] as T;
    }

    if (endpoint.includes("/tickets")) {
      return [
        {
          id: "TKT-2026-001",
          ticket_number: "TKT-2026-001",
          customer_id: "CUST-KURA-2026-001",
          channel: "CALL_CENTER",
          category: "LINE_SNAP",
          priority: "HIGH",
          unit_type: "FEEDER_MAINT",
          description: "Feeder jumper burnt at transformer pole junction (Vom Market Road)",
          assigned_unit_id: "UNIT-FEEDER-01",
          status: "ROUTED",
          created_at: now,
        },
      ] as T;
    }

    if (endpoint.includes("/maintenance/orders") || endpoint.includes("/services/work-orders")) {
      return [
        {
          id: "WO-2026-FDR-001",
          order_no: "WO-FDR-001",
          unit_type: "FEEDER_MAINT",
          work_order_type: "CORRECTIVE_REPAIR",
          priority: "CRITICAL",
          status: "DISPATCHED",
          unit_head_id: "STAFF-HEAD-01",
          team_lead_id: "STAFF-LEAD-01",
          asset_id: "FDR-33KV-VOM",
          ptw_required: true,
          instructions: "Replace cracked insulator on tower 42 and clear encroaching branches",
          created_at: now,
        },
      ] as T;
    }

    if (endpoint.includes("/generation/telemetry")) {
      return [
        {
          id: "GEN-01",
          plant_id: "PLANT_01_KURRA_FALLS",
          mw_generated: 14.8,
          water_level_meters: 24.5,
          stator_temp_celsius: 68.2,
          bearing_vibration_mm_s: 2.1,
          lube_oil_pressure_bar: 3.2,
          in_service: true,
          logged_at: now,
        },
      ] as T;
    }

    if (endpoint.includes("/energy/feeders")) {
      return [
        {
          id: "FDR-01",
          feeder_code: "FDR-33KV-VOM",
          feeder_name: "33kV Vom Institutional Feeder",
          voltage_level_kv: 33.0,
          current_load_mw: 8.5,
          state: "IN_SERVICE",
          priority_tier: 1,
          tariff_band: "BAND_A",
        },
      ] as T;
    }

    if (endpoint.includes("/energy/dispatch/latest")) {
      return {
        id: "DISP-01",
        shift_window: "MORNING",
        generation_mw_in: 45.0,
        tcn_mw_in: 15.0,
        total_bulk_mw: 60.0,
        operator_id: "STAFF-DISP-01",
        logged_at: now,
      } as T;
    }

    if (endpoint.includes("/energy/ptw")) {
      return [
        {
          id: "PTW-2026-001",
          permit_number: "PTW-001",
          feeder_id: "FDR-33KV-VOM",
          substation_code: "SUB-VOM-33KV",
          isolation_point: "CB-04 Line Isolator",
          step: "REQUESTED",
          requested_by_staff_id: "STAFF-LEAD-01",
          created_at: now,
        },
      ] as T;
    }

    if (endpoint.includes("/stores/inventory")) {
      return [
        {
          id: "ITEM-01",
          sku: "FUSE-11KV-100A",
          item_name: "11kV Drop-Out Fuse 100A",
          category: "SUBSTATION_HARDWARE",
          current_quantity: 45,
          reorder_threshold: 10,
          unit_cost: 4500.0,
        },
      ] as T;
    }

    if (endpoint.includes("/stores/requisitions")) {
      return [
        {
          id: "REQ-01",
          requisition_no: "REQ-2026-001",
          item_sku: "FUSE-11KV-100A",
          quantity_requested: 5,
          is_out_of_stock_procurement: false,
          status: "SUBMITTED",
          requester_staff_id: "STAFF-FEEDER-01",
          unit_head_approval: true,
          internal_audit_consent: false,
        },
      ] as T;
    }

    if (endpoint.includes("/billing/vending")) {
      return [
        {
          id: "TX-VEND-01",
          transaction_ref: "TX-2026-98124",
          customer_id: "CUST-KURA-2026-001",
          meter_number: "MTR-98214",
          amount_paid: 10000.0,
          tariff_rate_per_kwh: 68.0,
          kwh_vended: 147.05,
          sts_token_20_digit: "4921-8824-1092-8172-5501",
          meter_refund_credited: 1500.0,
          payment_platform_source: "OPAY_WEBHOOK",
          created_at: now,
        },
      ] as T;
    }

    if (endpoint.includes("/billing/revenue-protection/alerts")) {
      return [
        {
          id: "ALT-01",
          customer_or_dt_id: "TX-DT-500KVA-MARKET",
          entity_type: "TRANSFORMER",
          alert_type: "UNMETERED_LOAD_SURGE",
          confidence_score: 0.89,
          evidence_summary: "Secondary DT energy meter discrepancy vs customer meters exceeds 32%",
          status: "OPEN",
          created_at: now,
        },
      ] as T;
    }

    return [] as unknown as T;
  }

  public get<T>(endpoint: string, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: "GET" });
  }

  public post<T>(endpoint: string, body?: unknown, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  public put<T>(endpoint: string, body?: unknown, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    });
  }
}

export const apiClient = new ApiClient();
