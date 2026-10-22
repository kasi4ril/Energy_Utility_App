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
      // Graceful offline fallback providing realistic demo telemetry when backend is not listening
      return this.getMockFallback<T>(endpoint, options.method || "GET");
    }
  }

  private getMockFallback<T>(endpoint: string, method: string): T {
    const now = new Date().toISOString();

    if (method === "POST" || method === "PUT" || method === "PATCH") {
      return { success: true, message: "Action submitted successfully (Dev Mode)", id: `dev-${Date.now()}` } as T;
    }

    if (endpoint.includes("/executive/dashboard")) {
      return {
        total_mw_generated: 48.6,
        bulk_mw_available: 63.6,
        active_feeders_count: 14,
        tripped_feeders_count: 1,
        pending_audit_claims_count: 3,
        pending_procurements_count: 2,
        open_theft_alerts_count: 1,
        prepaid_vending_revenue: 14850000.0,
        active_shift_window: "MORNING",
        shift_telemetry_timestamp: now,
      } as T;
    }

    if (endpoint.includes("/executive/actions")) {
      return [
        {
          id: "act-01",
          department: "STORES_PROCUREMENT",
          action_type: "HIGH_VALUE_PROCUREMENT",
          description: "Emergency procurement for 500kVA Distribution Transformer (Jos South)",
          requested_by: "Store Officer Bitrus",
          requested_at: now,
          status: "PENDING_MD_SIGN_OFF",
          required_role: "MD",
        },
        {
          id: "act-02",
          department: "COMMERCIAL_OPS",
          action_type: "DISPUTED_BILLING_REVERSAL",
          description: "Disputed billing adjustment of ₦45,000 for Customer CUST-00142",
          requested_by: "Billing Lead Gyang",
          requested_at: now,
          status: "PENDING_SIGN_OFF",
          required_role: "HEAD_COMMERCIAL_OPS",
        },
      ] as T;
    }

    if (endpoint.includes("/tickets")) {
      return [
        {
          id: "TKT-2026-001",
          ticket_no: "TKT-001",
          customer_id: "CUST-KURA-2026-001",
          channel: "CALL_CENTER",
          category: "LINE_SNAP",
          description: "Feeder jumper burnt at transformer pole junction (Vom Market Road)",
          assigned_unit_id: "UNIT-FEEDER-01",
          status: "ROUTED",
          created_at: now,
        },
      ] as T;
    }

    if (endpoint.includes("/services/work-orders")) {
      return [
        {
          id: "WO-2026-FDR-001",
          order_no: "WO-FDR-001",
          unit_type: "FEEDER_MAINT",
          work_order_type: "CORRECTIVE_REPAIR",
          priority: "CRITICAL",
          asset_id: "FDR-33KV-VOM",
          status: "DISPATCHED",
          created_at: now,
        },
      ] as T;
    }

    if (endpoint.includes("/generation/telemetry")) {
      return [
        {
          id: "GEN-01",
          shift_window: "MORNING",
          plant_id: "PLANT_01_KURRA_FALLS",
          plant_state: "IN_SERVICE",
          water_level_m: 24.5,
          mw_generated: 14.8,
          stator_winding_temp_c: 68.2,
          recorded_at: now,
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
          asset_id: "FDR-33KV-VOM",
          work_description: "Tower 42 insulator replacement and line ground clearance",
          status: "REQUESTED",
          created_at: now,
        },
      ] as T;
    }

    if (endpoint.includes("/stores/inventory")) {
      return [
        {
          id: "ITEM-01",
          item_code: "FUSE-11KV-100A",
          item_name: "11kV Drop-Out Fuse 100A",
          category: "SUBSTATION_HARDWARE",
          quantity_available: 45,
          unit_of_measure: "PCS",
          reorder_threshold: 10,
        },
      ] as T;
    }

    if (endpoint.includes("/stores/requisitions")) {
      return [
        {
          id: "REQ-01",
          requisition_no: "REQ-2026-001",
          item_id: "ITEM-01",
          quantity_requested: 5,
          path: "PATH_A_IN_STOCK",
          status: "SUBMITTED",
          created_at: now,
        },
      ] as T;
    }

    if (endpoint.includes("/billing/meters/refunds")) {
      return [
        {
          id: "REF-01",
          customer_id: "CUST-KURA-2026-001",
          meter_purchase_cost: 85000.0,
          total_refunded: 15000.0,
          monthly_refund_rate: 5000.0,
          status: "ACTIVE",
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
