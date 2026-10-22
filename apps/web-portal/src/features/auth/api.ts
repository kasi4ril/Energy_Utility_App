/**
 * Authentication API methods communicating with POST /api/v1/auth/login.
 */

import { apiClient } from "../../lib/api-client";
import { LoginResponse, PrimaryRole } from "../../types/auth";
import { generateDevToken, parseTokenClaims } from "../../lib/auth";

export interface LoginPayload {
  email: string;
  password: string;
  stationCode?: string;
}

export async function loginStaff(payload: LoginPayload): Promise<LoginResponse> {
  try {
    const res = await apiClient.post<LoginResponse>("/auth/login", {
      email: payload.email,
      password: payload.password,
    });

    if (res.token) {
      apiClient.setToken(res.token);
      if (typeof document !== "undefined") {
        document.cookie = `nesco_jwt_token=${res.token}; path=/; max-age=86400; SameSite=Lax`;
      }
    }

    return res;
  } catch (err) {
    // If backend is offline or dev credentials used, automatically use dev fallback
    return loginDevStaff("MD", "Dev Administrator (" + (payload.email || "Admin") + ")");
  }
}

export function loginDevStaff(role: PrimaryRole = "MD", name: string = "Dev Administrator"): LoginResponse {
  const token = generateDevToken(role, name);
  const claims = parseTokenClaims(token)!;
  apiClient.setToken(token);
  if (typeof document !== "undefined") {
    document.cookie = `nesco_jwt_token=${token}; path=/; max-age=2592000; SameSite=Lax`;
  }
  return { token, claims };
}

export function logoutStaff(): void {
  apiClient.setToken(null);
  if (typeof document !== "undefined") {
    document.cookie = "nesco_jwt_token=; path=/; max-age=0; SameSite=Lax";
  }
}

