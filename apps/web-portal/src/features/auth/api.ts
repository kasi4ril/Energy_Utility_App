/**
 * Authentication API methods communicating with POST /api/v1/auth/login.
 */

import { apiClient } from "../../lib/api-client";
import { LoginResponse } from "../../types/auth";

export interface LoginPayload {
  email: string;
  password: string;
  stationCode?: string;
}

export async function loginStaff(payload: LoginPayload): Promise<LoginResponse> {
  const res = await apiClient.post<LoginResponse>("/auth/login", {
    email: payload.email,
    password: payload.password,
  });

  if (res.token) {
    apiClient.setToken(res.token);
    // Also store cookie for Next.js SSR middleware inspection
    if (typeof document !== "undefined") {
      document.cookie = `nesco_jwt_token=${res.token}; path=/; max-age=900; SameSite=Strict; Secure`;
    }
  }

  return res;
}

export function logoutStaff(): void {
  apiClient.setToken(null);
  if (typeof document !== "undefined") {
    document.cookie = "nesco_jwt_token=; path=/; max-age=0";
  }
}

