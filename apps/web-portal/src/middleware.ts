/**
 * Next.js Edge Middleware enforcing authentication guards and role-based auto-routing.
 * Directs staff seamlessly to their authorized departmental portal.
 */

import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { getLandingPageForRole, parseTokenClaims } from "./lib/auth";
import { PrimaryRole } from "./types/auth";

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // 1. Allow public assets and login route
  if (
    pathname.startsWith("/_next") ||
    pathname.startsWith("/api") ||
    pathname.startsWith("/static") ||
    pathname === "/login" ||
    pathname === "/favicon.ico"
  ) {
    return NextResponse.next();
  }

  // 2. Extract JWT token from cookie or Authorization header
  const token =
    request.cookies.get("nesco_jwt_token")?.value ||
    request.headers.get("Authorization")?.replace("Bearer ", "");

  if (!token) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("from", pathname);
    return NextResponse.redirect(loginUrl);
  }

  // 3. Inspect claims and expiration
  const claims = parseTokenClaims(token);
  if (!claims || (claims.exp && claims.exp * 1000 < Date.now())) {
    const loginUrl = new URL("/login", request.url);
    const response = NextResponse.redirect(loginUrl);
    response.cookies.delete("nesco_jwt_token");
    return response;
  }

  // 4. Root path auto-redirection to the role's assigned landing page
  if (pathname === "/") {
    const destination = getLandingPageForRole(claims.primary_role as PrimaryRole);
    return NextResponse.redirect(new URL(destination, request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico).*)"],
};

