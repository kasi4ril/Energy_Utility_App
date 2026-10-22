/**
 * Next.js Page: Single Unified Login Portal (/login).
 */

import React from "react";
import Head from "next/head";
import { LoginPage } from "../features/auth/LoginPage";

export default function LoginRoute() {
  return (
    <>
      <Head>
        <title>NESCO Login &bull; Unified Enterprise Command</title>
      </Head>
      <LoginPage />
    </>
  );
}

