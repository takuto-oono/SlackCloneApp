import { AppProps } from "next/app";
import { ReactElement, ReactNode } from "react";
import { CookiesProvider } from "react-cookie";
import 'src/styles/globals.css'
import Layout from "./common/layout";
import { NextPage } from "next/types";
import { RecoilRoot } from "recoil";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <CookiesProvider>
      <Layout >
        <Component {...pageProps} />
      </Layout>
    </CookiesProvider>
  )
}
