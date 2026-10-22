/**
 * Global Next.js Application Wrapper (_app.tsx).
 * Injects typography, global styles, and enterprise dark mode tokens.
 */

import React from "react";
import type { AppProps } from "next/app";
import Head from "next/head";
import Script from "next/script";
import "../styles/globals.css";

export default function NescoApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <Head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Script src="https://cdn.tailwindcss.com" strategy="beforeInteractive" />
      <Script
        id="tailwind-config"
        strategy="beforeInteractive"
        dangerouslySetInnerHTML={{
          __html: `
            tailwind.config = {
              darkMode: 'class',
              theme: {
                extend: {
                  colors: {
                    slate: {
                      850: '#172033',
                      950: '#0a0f1d',
                    }
                  }
                }
              }
            }
          `,
        }}
      />
      <Component {...pageProps} />
    </>
  );
}

